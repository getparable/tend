// Copyright (c) 2026 Parable, LLC
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	b, err := exec.Command("xcrun", "--find", "swiftc").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	plugin := filepath.Join(filepath.Dir(filepath.Dir(strings.TrimSpace(string(b)))), "lib/swift/host/plugins/testing/libTestingMacros.dylib")
	args := []string{"test", "--package-path", "macos"}
	// Swift 6.4 CLT ships the macro in a nested directory that SwiftPM's new
	// default build system does not automatically discover. Load that exact
	// toolchain's plugin when present; full Xcode can use its normal discovery.
	if _, err = os.Stat(plugin); err == nil {
		args = append(args, "-Xswiftc", "-load-plugin-library", "-Xswiftc", plugin)
	}
	fmt.Printf("$ swift %q\n", args)
	c := exec.Command("swift", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err = c.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
