// Copyright (c) 2026 Parable, LLC
package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/getparable/tend/internal/ipc"
	"github.com/getparable/tend/internal/protocol"
)

func TestLifecycle(t *testing.T) {
	// Canonical short root: Darwin sockaddr_un has only 104 path bytes.
	base, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root, err := os.MkdirTemp(base, "td-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)
	path := filepath.Join(root, "s")
	if len(path) > 103 {
		t.Skip("TMPDIR too long for Darwin Unix socket; run with a short isolated TMPDIR")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, path) }()
	deadline := time.Now().Add(3 * time.Second)
	for ipc.ValidateSocket(path) != nil {
		select {
		case err := <-done:
			t.Fatal(err)
		default:
		}
		if time.Now().After(deadline) {
			t.Fatal("socket timeout")
		}
		time.Sleep(10 * time.Millisecond)
	}
	call := func(method, fixture string) protocol.Response {
		return ipc.Call(path, protocol.Request{Version: 1, ID: "test", Method: method, Fixture: fixture})
	}
	if r := call("fixture.start", "complete"); r.Error != nil {
		t.Fatal(r.Error)
	}
	if r := call("fixture.start", "complete"); r.Error == nil || r.Error.Code != "already_running" {
		t.Fatal(r)
	}
	wait := func(state string) {
		t.Helper()
		deadline := time.Now().Add(6 * time.Second)
		for {
			r := call("service.status", "")
			if r.Error != nil {
				t.Fatal(r.Error)
			}
			if r.Result.Run.State == state {
				return
			}
			if time.Now().After(deadline) {
				t.Fatalf("expected %s: %+v", state, r.Result.Run)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	wait("succeeded") // New socket connections see the same in-memory result.
	call("fixture.start", "cancel")
	time.Sleep(100 * time.Millisecond)
	call("fixture.cancel", "")
	wait("cancelled")
	call("fixture.start", "cancel")
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(6 * time.Second):
		t.Fatal("unbounded teardown")
	}
}
