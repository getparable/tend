// Copyright (c) 2026 Parable, LLC
package ipc

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/getparable/tend/internal/protocol"
)

func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "TendLifecycle-dev", "service.sock")
}
func ValidatePath(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path || len(path) > 103 {
		return errors.New("socket must be an absolute clean path of at most 103 bytes")
	}
	return validateDirectory(filepath.Dir(path))
}
func validateDirectory(root string) error {
	for p := root; ; p = filepath.Dir(p) {
		st, err := os.Lstat(p)
		if err != nil {
			return err
		}
		if st.Mode()&os.ModeSymlink != 0 || !st.IsDir() {
			return errors.New("socket ancestors must be real directories")
		}
		if p == root && (st.Mode().Perm() != 0700 || st.Sys().(*syscall.Stat_t).Uid != uint32(os.Geteuid())) {
			return errors.New("socket directory must be owned by this user with mode 0700")
		}
		if p == "/" {
			break
		}
	}
	return nil
}
func ValidateSocket(path string) error {
	if err := ValidatePath(path); err != nil {
		return err
	}
	return validateSocketEntry(path)
}
func validateSocketEntry(path string) error {
	st, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if st.Mode()&os.ModeSocket == 0 || st.Mode().Perm() != 0600 || st.Sys().(*syscall.Stat_t).Uid != uint32(os.Geteuid()) {
		return errors.New("socket must be owned by this user with mode 0600")
	}
	return nil
}
func Call(path string, req protocol.Request) protocol.Response {
	fail := func(code string, err error) protocol.Response { return protocol.Failure(req.ID, code, err.Error()) }
	if err := ValidateSocket(path); err != nil {
		return fail("service_unavailable", err)
	}
	c, err := net.DialTimeout("unix", path, time.Second)
	if err != nil {
		return fail("service_unavailable", err)
	}
	defer c.Close()
	if err = checkPeer(c.(*net.UnixConn)); err != nil {
		return fail("unsafe_peer", err)
	}
	_ = c.SetDeadline(time.Now().Add(5 * time.Second))
	if err = protocol.Write(c, req); err != nil {
		return fail("invalid_request", err)
	}
	var resp protocol.Response
	if err = protocol.Read(c, &resp); err != nil {
		return fail("invalid_response", err)
	}
	if resp.Version != protocol.Version {
		return fail("incompatible_version", fmt.Errorf("helper protocol is %d", resp.Version))
	}
	if resp.ID != req.ID || (resp.Result == nil) == (resp.Error == nil) {
		return fail("invalid_response", errors.New("invalid response identity or result"))
	}
	return resp
}
func CheckPeer(c *net.UnixConn) error { return checkPeer(c) }
