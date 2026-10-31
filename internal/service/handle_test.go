// Copyright (c) 2026 Parable, LLC
package service

import (
	"testing"
	"time"

	"github.com/getparable/tend/internal/protocol"
)

func TestInMemoryAdmissionAndCancellation(t *testing.T) {
	s := &Service{}
	call := func(method, fixture string) protocol.Response {
		return s.Handle(protocol.Request{Version: 1, ID: "direct", Method: method, Fixture: fixture})
	}
	if r := call("fixture.start", "cancel"); r.Error != nil {
		t.Fatal(r.Error)
	}
	if r := call("fixture.start", "cancel"); r.Error == nil || r.Error.Code != "already_running" {
		t.Fatal(r)
	}
	time.Sleep(100 * time.Millisecond)
	call("fixture.cancel", "")
	deadline := time.Now().Add(4 * time.Second)
	for {
		r := call("service.status", "")
		if r.Result.Run.State == "cancelled" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("cancellation did not finish")
		}
		time.Sleep(20 * time.Millisecond)
	}
	s.mu.Lock()
	s.closing = true
	s.mu.Unlock()
	if r := call("fixture.start", "complete"); r.Error == nil || r.Error.Code != "service_unavailable" {
		t.Fatal(r)
	}
}
