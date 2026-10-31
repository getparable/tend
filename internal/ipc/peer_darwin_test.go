// Copyright (c) 2026 Parable, LLC
package ipc

import (
	"net"
	"os"
	"syscall"
	"testing"
)

func TestPeerIdentity(t *testing.T) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, fd := range fds {
		f := os.NewFile(uintptr(fd), "socketpair")
		c, err := net.FileConn(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		if err = CheckPeer(c.(*net.UnixConn)); err != nil {
			t.Fatal(err)
		}
		c.Close()
	}
}
