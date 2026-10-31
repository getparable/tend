// Copyright (c) 2026 Parable, LLC
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/getparable/tend/internal/ipc"
	"github.com/getparable/tend/internal/protocol"
)

func main() {
	if err := verify(); err != nil {
		fmt.Fprintln(os.Stderr, "BLOCKED/FAILED:", err)
		os.Exit(1)
	}
}
func command(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, args[0], args[1:]...)
	b, err := c.CombinedOutput()
	fmt.Printf("$ %q\n%s", args, b)
	return string(b), err
}
func verify() (retErr error) {
	mode := flag.String("mode", "attended", "launchd or attended")
	flag.Parse()
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	app := filepath.Join(home, "Applications", "Tend Lifecycle.app")
	exe := filepath.Join(app, "Contents/MacOS/Tend")
	if _, err = os.Stat(exe); err != nil {
		return fmt.Errorf("run make app and make install-dev first: %w", err)
	}
	if *mode == "attended" {
		return attended(app, exe)
	}
	if *mode != "launchd" {
		return fmt.Errorf("unknown verification mode")
	}
	domain := fmt.Sprintf("gui/%d", os.Getuid())
	probe := exec.Command("/bin/launchctl", "print", domain)
	probe.Stdout = io.Discard
	probe.Stderr = io.Discard
	if err = probe.Run(); err != nil {
		return fmt.Errorf("signed-in GUI bootstrap domain required: %w", err)
	}
	path := filepath.Join(home, "Library", "Application Support", "TendLifecycle-proof", "service.sock")
	root := filepath.Dir(path)
	// This is a dedicated installed runtime location for the test-label helper.
	// Never share the real SMAppService socket, even if registration is disabled.
	if _, err = os.Lstat(root); !os.IsNotExist(err) {
		return fmt.Errorf("refusing existing runtime state %s; unregister and clean only Tend lifecycle state first", root)
	}
	if err = os.Mkdir(root, 0700); err != nil {
		return err
	}
	defer os.RemoveAll(root)
	if err = ipc.ValidatePath(path); err != nil {
		return err
	}
	label := fmt.Sprintf("io.parable.tend.lifecycle.test.%d.%d", os.Getpid(), time.Now().UnixNano())
	job := domain + "/" + label
	f, err := os.CreateTemp("", "tend-launchd-*.plist")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	// JSON is accepted by plutil; XML conversion avoids manual plist escaping.
	plist := map[string]any{"Label": label, "ProgramArguments": []string{filepath.Join(app, "Contents/Helpers/tend"), "serve", "--socket", path}, "RunAtLoad": true, "KeepAlive": true, "ExitTimeOut": 5, "ThrottleInterval": 1}
	err = json.NewEncoder(f).Encode(plist)
	f.Close()
	if err != nil {
		return err
	}
	if _, err = command("/usr/bin/plutil", "-convert", "xml1", f.Name()); err != nil {
		return err
	}
	if _, err = command("/bin/launchctl", "bootstrap", domain, f.Name()); err != nil {
		return err
	}
	loaded := true
	defer func() {
		if loaded {
			if _, err := command("/bin/launchctl", "bootout", job); err != nil && retErr == nil {
				retErr = fmt.Errorf("test job cleanup failed: %w", err)
			}
		}
	}()
	call := func(method, fixture string) protocol.Response {
		r := ipc.Call(path, protocol.Request{Version: 1, ID: "lifecycle", Method: method, Fixture: fixture})
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
		return r
	}
	wait := func(test func(*protocol.Result) bool) error {
		deadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(deadline) {
			r := call("service.status", "")
			if r.Result != nil && test(r.Result) {
				return nil
			}
			time.Sleep(250 * time.Millisecond)
		}
		return fmt.Errorf("timed out waiting for helper state")
	}
	if err = wait(func(r *protocol.Result) bool { return r.PID > 0 }); err != nil {
		return err
	}
	r := call("fixture.start", "complete")
	if r.Error != nil {
		return fmt.Errorf("start: %+v", r.Error)
	}
	if err = wait(func(r *protocol.Result) bool {
		return r.Run != nil && r.Run.State == "succeeded" && strings.Contains(r.Run.Output, "completed")
	}); err != nil {
		return err
	}
	pid := call("service.status", "").Result.PID
	// Request launchd to signal only this unique test job, never a saved PID.
	if _, err = command("/bin/launchctl", "kill", "SIGTERM", job); err != nil {
		return err
	}
	if err = wait(func(r *protocol.Result) bool { return r.PID != pid && r.Run == nil }); err != nil {
		return err
	}
	if r = call("fixture.start", "cancel"); r.Error != nil {
		return fmt.Errorf("cancel start: %+v", r.Error)
	}
	time.Sleep(200 * time.Millisecond)
	call("fixture.cancel", "")
	if err = wait(func(r *protocol.Result) bool { return r.Run != nil && r.Run.State == "cancelled" }); err != nil {
		return err
	}
	call("fixture.start", "cancel")
	time.Sleep(200 * time.Millisecond)
	if _, err = command("/bin/launchctl", "bootout", job); err != nil {
		return err
	}
	loaded = false
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		if _, err = os.Lstat(path); os.IsNotExist(err) {
			fmt.Println("PASS: test-label start, reconnect, KeepAlive restart, cancellation and socket teardown. SMAppService/UI approval remains unverified.")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("socket teardown did not complete")
}
