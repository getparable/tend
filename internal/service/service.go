// Copyright (c) 2026 Parable, LLC
package service

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/getparable/tend/internal/ipc"
	"github.com/getparable/tend/internal/protocol"
	"github.com/getparable/tend/internal/runner"
)

type Service struct {
	mu      sync.Mutex
	child   *runner.Child
	run     *protocol.Run
	closing bool
}

func (s *Service) Handle(req protocol.Request) protocol.Response {
	if e := protocol.Validate(req); e != nil {
		return protocol.Failure(req.ID, e.Code, e.Message)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closing {
		return protocol.Failure(req.ID, "service_unavailable", "Helper is stopping")
	}
	switch req.Method {
	case "fixture.start":
		if s.child != nil {
			return protocol.Failure(req.ID, "already_running", "Fixture is already running")
		}
		c, err := runner.Start(req.Fixture)
		if err != nil {
			return protocol.Failure(req.ID, "spawn_failed", err.Error())
		}
		s.child = c
		s.run = &protocol.Run{State: "running", Output: ""}
		go func() {
			<-c.Done()
			s.mu.Lock()
			defer s.mu.Unlock()
			code := c.ExitCode()
			state := "succeeded"
			if code != 0 {
				state = "failed"
			}
			if s.run.State == "stopping" {
				state = "cancelled"
			}
			s.run = &protocol.Run{State: state, Output: c.Output(), ExitCode: &code}
			s.child = nil
		}()
	case "fixture.cancel":
		if s.child != nil && s.run.State != "stopping" {
			s.run.State = "stopping"
			go s.child.Cancel()
		}
	}
	var run *protocol.Run
	if s.run != nil {
		copy := *s.run
		run = &copy
	}
	return protocol.Response{Version: protocol.Version, ID: req.ID, Result: &protocol.Result{PID: os.Getpid(), Run: run}}
}
func Serve(ctx context.Context, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if err := ipc.ValidatePath(path); err != nil {
		return err
	}
	// O_NOFOLLOW and flock prevent replacing a running helper's socket.
	fd, err := syscall.Open(path+".lock", syscall.O_CREAT|syscall.O_RDWR|syscall.O_NOFOLLOW, 0600)
	if err != nil {
		return err
	}
	lock := os.NewFile(uintptr(fd), path+".lock")
	defer lock.Close()
	if err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return errors.New("helper already owns socket")
	}
	if _, err = os.Lstat(path); err == nil {
		if err = ipc.ValidateSocket(path); err != nil {
			return err
		}
		if err = os.Remove(path); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	// Set process-wide umask once, before accepting requests or starting children.
	syscall.Umask(0077)
	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		return err
	}
	defer l.Close()
	if err = os.Chmod(path, 0600); err != nil {
		return err
	}
	s := &Service{}
	stopped := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			l.Close()
		case <-stopped:
		}
	}()
	defer close(stopped)
	var wg sync.WaitGroup
	defer func() {
		s.mu.Lock()
		s.closing = true
		c := s.child
		s.mu.Unlock()
		if c != nil {
			c.Cancel()
		}
		wg.Wait()
	}()
	// Bound simultaneous connections as well as individual frame size/deadline.
	slots := make(chan struct{}, 16)
	for {
		c, err := l.AcceptUnix()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		select {
		case slots <- struct{}{}:
		default:
			c.Close()
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			defer c.Close()
			_ = c.SetDeadline(time.Now().Add(5 * time.Second))
			if ipc.CheckPeer(c) != nil {
				return
			}
			var req protocol.Request
			if err := protocol.Read(c, &req); err != nil {
				_ = protocol.Write(c, protocol.Failure(req.ID, "invalid_request", err.Error()))
				return
			}
			_ = protocol.Write(c, s.Handle(req))
		}()
	}
}
