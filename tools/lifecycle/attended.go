// Copyright (c) 2026 Parable, LLC
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/getparable/tend/internal/ipc"
	"github.com/getparable/tend/internal/protocol"
)

type uiProcess struct {
	cmd  *exec.Cmd
	done chan error
}

func launchUI(exe, probe string) (*uiProcess, error) {
	c := exec.Command(exe)
	c.Env = append(os.Environ(), "TEND_LIFECYCLE_PROBE="+probe)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return nil, err
	}
	u := &uiProcess{c, make(chan error, 1)}
	go func() { u.done <- c.Wait() }()
	fmt.Printf("Launched owned UI PID %d\n", c.Process.Pid)
	return u, nil
}
func (u *uiProcess) quit(exe string) error {
	if _, err := command(exe, "--quit-owned-ui", fmt.Sprint(u.cmd.Process.Pid)); err != nil {
		u.cmd.Process.Kill()
		return err
	}
	select {
	case err := <-u.done:
		fmt.Println("Owned UI exited:", err)
		return err
	case <-time.After(5 * time.Second):
		u.cmd.Process.Kill()
		return fmt.Errorf("UI did not quit within five seconds")
	}
}
func awaitUI(probe, state string) error {
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		b, err := os.ReadFile(probe)
		var snapshot struct {
			Connection string `json:"connection"`
			State      string `json:"run_state"`
		}
		if err == nil && json.Unmarshal(b, &snapshot) == nil && snapshot.Connection == "Available" && (state == "" || snapshot.State == state) {
			fmt.Printf("Installed UI model observation: %s\n", b)
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("no installed UI observation for state %q; GUI session required", state)
}
func attended(app, exe string) (retErr error) {
	domain := fmt.Sprintf("gui/%d", os.Getuid())
	if err := exec.Command("/bin/launchctl", "print", domain).Run(); err != nil {
		return fmt.Errorf("signed-in GUI session required: %w", err)
	}
	out, err := command(exe, "--helper", "status")
	if err != nil {
		return fmt.Errorf("helper requires approval or is disabled; review System Settings: %w", err)
	}
	owned := false
	if !strings.Contains(out, "status: 1 ") {
		if _, err = command(exe, "--helper", "register"); err != nil {
			return fmt.Errorf("registration/approval required: %w", err)
		}
		owned = true
	}
	path := ipc.DefaultPath()
	call := func(method, fixture string) (protocol.Response, error) {
		req := protocol.Request{Version: 1, ID: "installed-lifecycle", Method: method, Fixture: fixture}
		var input bytes.Buffer
		if err := protocol.Write(&input, req); err != nil {
			return protocol.Response{}, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		c := exec.CommandContext(ctx, filepath.Join(app, "Contents/Helpers/tend"), "request")
		c.Stdin = &input
		b, err := c.Output()
		fmt.Printf("CLI %s: %s", method, b)
		var resp protocol.Response
		if decode := protocol.Read(bytes.NewReader(b), &resp); decode != nil {
			return resp, decode
		}
		if resp.Error != nil {
			return resp, fmt.Errorf("%s: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp, err
	}
	wait := func(test func(*protocol.Result) bool) error {
		deadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(deadline) {
			r, err := call("service.status", "")
			if err == nil && r.Result != nil && test(r.Result) {
				return nil
			}
			time.Sleep(250 * time.Millisecond)
		}
		return fmt.Errorf("installed helper state timed out")
	}
	var ui *uiProcess
	defer func() {
		if !owned {
			return
		}
		if ui != nil {
			if err := ui.quit(exe); err != nil && retErr == nil {
				retErr = err
			}
		}
		if _, err := command(exe, "--helper", "unregister"); err != nil {
			if retErr == nil {
				retErr = err
			}
			return
		}
		deadline := time.Now().Add(8 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := os.Lstat(path); os.IsNotExist(err) {
				// Remove only the known lock and an empty proof directory.
				os.Remove(path + ".lock")
				os.Remove(filepath.Dir(path))
				fmt.Println("Cleanup: unregistered, socket removed, owned UI and fixture stopped.")
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		if retErr == nil {
			retErr = fmt.Errorf("unregister did not remove socket within eight seconds")
		}
	}()
	if err = wait(func(r *protocol.Result) bool { return r.PID > 0 }); err != nil {
		_, _ = command("/bin/launchctl", "print", domain+"/io.parable.tend.lifecycle.agent")
		return fmt.Errorf("%w; Enabled describes registration, not reachability. Inspect the job state, executable, socket and scoped launch diagnostics before attributing the failure to signing", err)
	}
	first, err := call("service.status", "")
	if err != nil {
		return err
	}
	if first.Result.Run != nil && (first.Result.Run.State == "running" || first.Result.Run.State == "stopping") {
		return fmt.Errorf("existing fixture is active; preserve it and retry after it finishes")
	}
	owned = true
	f, err := os.CreateTemp("", "tend-ui-observation-*.json")
	if err != nil {
		return err
	}
	probe := f.Name()
	f.Close()
	os.Remove(probe)
	defer os.Remove(probe)
	ui, err = launchUI(exe, probe)
	if err != nil {
		return err
	}
	if err = awaitUI(probe, ""); err != nil {
		return err
	}
	if _, err = call("fixture.start", "complete"); err != nil {
		return err
	}
	if err = awaitUI(probe, "running"); err != nil {
		return err
	}
	err = ui.quit(exe)
	ui = nil
	if err != nil {
		return err
	}
	r, err := call("service.status", "")
	if err != nil {
		return err
	}
	if r.Result.Run == nil || r.Result.Run.State != "running" {
		return fmt.Errorf("fixture did not remain running after UI exit; rerun on an idle machine")
	}
	if err = wait(func(r *protocol.Result) bool {
		return r.Run != nil && r.Run.State == "succeeded" && strings.Contains(r.Run.Output, "completed")
	}); err != nil {
		return err
	}
	os.Remove(probe)
	ui, err = launchUI(exe, probe)
	if err != nil {
		return err
	}
	if err = awaitUI(probe, "succeeded"); err != nil {
		return err
	}
	pid := r.Result.PID
	job := fmt.Sprintf("gui/%d/io.parable.tend.lifecycle.agent", os.Getuid())
	if _, err = command("/bin/launchctl", "kill", "SIGTERM", job); err != nil {
		return err
	}
	if err = wait(func(r *protocol.Result) bool { return r.PID != pid && r.Run == nil }); err != nil {
		return err
	}
	if _, err = call("fixture.start", "cancel"); err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	if _, err = call("fixture.cancel", ""); err != nil {
		return err
	}
	if err = wait(func(r *protocol.Result) bool { return r.Run != nil && r.Run.State == "cancelled" }); err != nil {
		return err
	}
	if _, err = call("fixture.start", "cancel"); err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	fmt.Println("PASS: installed SMAppService registration, UI quit/reconnect, in-memory result, idle KeepAlive restart and bounded cancellation. Unregister/active-child teardown follows.")
	return nil
}
