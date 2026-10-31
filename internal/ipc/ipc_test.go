// Copyright (c) 2026 Parable, LLC
package ipc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnsafePaths(t *testing.T) {
	if ValidatePath("/"+strings.Repeat("a", 104)) == nil {
		t.Fatal("accepted long path")
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(root, 0755); err != nil {
		t.Fatal(err)
	}
	if validateDirectory(root) == nil {
		t.Fatal("accepted public root")
	}
	os.Chmod(root, 0700)
	if err = validateDirectory(root); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "s")
	os.WriteFile(file, []byte("keep"), 0600)
	if validateSocketEntry(file) == nil {
		t.Fatal("accepted regular file")
	}
	link := filepath.Join(root, "link")
	os.Symlink(root, link)
	if validateDirectory(link) == nil {
		t.Fatal("accepted symlink")
	}
}
