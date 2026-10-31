// Copyright (c) 2026 Parable, LLC
package ipc

/*
#include <sys/types.h>
#include <unistd.h>
static int same_user(int fd) {
  uid_t uid; gid_t gid;
  return getpeereid(fd, &uid, &gid) == 0 && uid == geteuid();
}
*/
import "C"
import (
	"errors"
	"net"
)

func checkPeer(c *net.UnixConn) error {
	raw, err := c.SyscallConn()
	if err != nil {
		return err
	}
	ok := false
	if err = raw.Control(func(fd uintptr) { ok = C.same_user(C.int(fd)) != 0 }); err != nil {
		return err
	}
	if !ok {
		return errors.New("peer uid does not match effective uid")
	}
	return nil
}
