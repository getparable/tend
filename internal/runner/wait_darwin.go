// Copyright (c) 2026 Parable, LLC
package runner

/*
#include <sys/wait.h>
#include <errno.h>
static int await_exit(int pid) {
  siginfo_t info;
  int result;
  do { result = waitid(P_PID, pid, &info, WEXITED | WNOWAIT); }
  while (result == -1 && errno == EINTR);
  return result;
}
*/
import "C"

// Observe exit without reaping. Holding the zombie keeps the process group ID
// from being reused until the signal/Wait handoff is complete.
func awaitExit(pid int) bool { return C.await_exit(C.int(pid)) == 0 }
