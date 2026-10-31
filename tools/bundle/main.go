// Copyright (c) 2026 Parable, LLC
// Build tools deliberately use argv arrays, never shell command strings.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run(name string, args ...string) {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	must(c.Run())
}
func main() {
	root, err := os.Getwd()
	must(err)
	app := filepath.Join(root, "dist", "Tend Lifecycle.app")
	if len(os.Args) == 2 && os.Args[1] == "install" {
		home, err := os.UserHomeDir()
		must(err)
		target := filepath.Join(home, "Applications", "Tend Lifecycle.app")
		if _, err = os.Lstat(target); !os.IsNotExist(err) {
			must(fmt.Errorf("refusing to replace %s; unregister and remove the previous development bundle explicitly", target))
		}
		must(os.MkdirAll(filepath.Dir(target), 0755))
		run("/usr/bin/ditto", app, target)
		fmt.Println("Installed isolated development bundle:", target)
		return
	}
	// This tool owns only this generated bundle, never the installed target.
	must(os.RemoveAll(app))
	for _, dir := range []string{"MacOS", "Helpers", "Resources", "Library/LaunchAgents"} {
		must(os.MkdirAll(filepath.Join(app, "Contents", dir), 0755))
	}
	run("go", "build", "-trimpath", "-o", filepath.Join(app, "Contents/Helpers/tend"), "./cmd/tend")
	run("swift", "build", "--package-path", "macos", "-c", "release", "--arch", "arm64")
	c := exec.Command("swift", "build", "--package-path", "macos", "-c", "release", "--arch", "arm64", "--show-bin-path")
	out, err := c.Output()
	must(err)
	bin := string(out)
	for len(bin) > 0 && (bin[len(bin)-1] == '\n' || bin[len(bin)-1] == '\r') {
		bin = bin[:len(bin)-1]
	}
	copyFile := func(src, dst string, mode os.FileMode) {
		b, err := os.ReadFile(src)
		must(err)
		must(os.WriteFile(dst, b, mode))
	}
	copyFile(filepath.Join(bin, "Tend"), filepath.Join(app, "Contents/MacOS/Tend"), 0755)
	copyFile("macos/Resources/Info.plist", filepath.Join(app, "Contents/Info.plist"), 0644)
	copyFile("macos/Resources/io.parable.tend.lifecycle.agent.plist", filepath.Join(app, "Contents/Library/LaunchAgents/io.parable.tend.lifecycle.agent.plist"), 0644)
	copyFile("LICENSE", filepath.Join(app, "Contents/Resources/LICENSE"), 0644)
	identity := os.Getenv("CODE_SIGN_IDENTITY")
	if identity == "" {
		identity = "-"
	}
	run("/usr/bin/codesign", "--force", "--sign", identity, "--identifier", "io.parable.tend.lifecycle.helper", filepath.Join(app, "Contents/Helpers/tend"))
	run("/usr/bin/codesign", "--force", "--sign", identity, app)
	run("/usr/bin/codesign", "--verify", "--deep", "--strict", app)
	if identity == "-" {
		fmt.Println("Built ad-hoc signed development bundle (SMAppService may require trusted signing):", app)
	} else {
		fmt.Println("Built development bundle with the explicitly supplied signing identity:", app)
	}
}
