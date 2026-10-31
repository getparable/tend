// Copyright (c) 2026 Parable, LLC
package runner

import (
	"context"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestFixtureCompletionAndCancellation(t *testing.T) {
	c, err := Start("complete")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(c.Cancel)
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	if err = c.Wait(ctx); err != nil {
		t.Fatal(err)
	}
	if c.ExitCode() != 0 || !strings.Contains(c.Output(), "completed") {
		t.Fatalf("unexpected result: %d %q", c.ExitCode(), c.Output())
	}
	c, err = Start("cancel")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(c.Cancel)
	time.Sleep(100 * time.Millisecond)
	start := time.Now()
	c.Cancel()
	if time.Since(start) > 4*time.Second {
		t.Fatal("unbounded cancellation")
	}
	select {
	case <-c.Done():
	case <-time.After(time.Second):
		t.Fatal("child survived cancellation")
	}
	if c.ExitCode() == 0 {
		t.Fatal("cancelled child reported success")
	}
	c.Cancel() // Idempotent cancellation never signals a completed child.
	deadline := time.Now().Add(2 * time.Second)
	for syscall.Kill(-c.cmd.Process.Pid, 0) != syscall.ESRCH {
		if time.Now().After(deadline) {
			t.Fatal("owned process group survived cancellation")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
