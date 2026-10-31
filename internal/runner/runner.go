// Copyright (c) 2026 Parable, LLC
package runner

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/getparable/tend/tests/fixtures"
)

type limitedBuffer struct{ b bytes.Buffer }

func (b *limitedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if left := 4096 - b.b.Len(); left > 0 {
		b.b.Write(p[:min(left, n)])
	}
	return n, nil
}

type Child struct {
	cmd    *exec.Cmd
	done   chan struct{}
	stop   sync.Once
	mu     sync.Mutex
	exited bool
	output limitedBuffer
	Err    error
}

func Start(name string) (*Child, error) {
	data, err := fixtures.Scripts.ReadFile("scripts/" + name + ".sh")
	if err != nil {
		return nil, err
	}
	f, writer, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	if _, err = writer.Write(data); err != nil {
		f.Close()
		writer.Close()
		return nil, err
	}
	writer.Close()
	c := &Child{cmd: exec.Command("/bin/sh", "/dev/fd/3"), done: make(chan struct{})}
	c.cmd.ExtraFiles = []*os.File{f}
	c.cmd.Env = []string{"PATH=/usr/bin:/bin", "LANG=C"}
	c.cmd.Dir = "/"
	c.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.cmd.Stdout = &c.output
	c.cmd.Stderr = &c.output
	c.cmd.WaitDelay = time.Second
	if err := c.cmd.Start(); err != nil {
		f.Close()
		return nil, err
	}
	f.Close()
	go func() {
		observed := awaitExit(c.cmd.Process.Pid)
		c.mu.Lock()
		c.exited = true
		if observed {
			_ = syscall.Kill(-c.cmd.Process.Pid, syscall.SIGKILL)
		}
		c.mu.Unlock()
		c.Err = c.cmd.Wait()
		close(c.done)
	}()
	return c, nil
}
func (c *Child) Done() <-chan struct{} { return c.done }
func (c *Child) Output() string        { <-c.done; return c.output.b.String() }
func (c *Child) ExitCode() int         { <-c.done; return c.cmd.ProcessState.ExitCode() }

// The child is still unreaped while signals are sent: never signal a saved PID.
// This proof supports attached process groups, not deliberately detached jobs.
func (c *Child) Cancel() {
	c.stop.Do(func() {
		if !c.signal(syscall.SIGTERM) {
			return
		}
		t := time.NewTimer(time.Second)
		defer t.Stop()
		select {
		case <-c.done:
			return
		case <-t.C:
		}
		c.signal(syscall.SIGKILL)
		select {
		case <-c.done:
		case <-time.After(2 * time.Second):
		}
	})
}
func (c *Child) signal(sig syscall.Signal) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.exited {
		return false
	}
	_ = syscall.Kill(-c.cmd.Process.Pid, sig)
	return true
}
func (c *Child) Wait(ctx context.Context) error {
	select {
	case <-c.done:
		return c.Err
	case <-ctx.Done():
		return ctx.Err()
	}
}
