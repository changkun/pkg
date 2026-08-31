// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package errhandling_test

import (
	"os"
	"path/filepath"
	"testing"

	errhandling "changkun.de/x/pkg/misc/error-handling-abstraction"
)

// copiers runs every test against both implementations, so the abstraction is
// held to the same contract as the explicit version.
var copiers = map[string]func(src, dst string) error{
	"CopyFile":     errhandling.CopyFile,
	"SafeCopyFile": errhandling.SafeCopyFile,
}

func TestCopyContent(t *testing.T) {
	want := []byte("the quick brown fox\njumps over the lazy dog\n")

	for name, copy := range copiers {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "src")
			dst := filepath.Join(dir, "dst")
			if err := os.WriteFile(src, want, 0o600); err != nil {
				t.Fatalf("write source: %v", err)
			}

			if err := copy(src, dst); err != nil {
				t.Fatalf("copy: %v", err)
			}

			got, err := os.ReadFile(dst)
			if err != nil {
				t.Fatalf("read destination: %v", err)
			}
			if string(got) != string(want) {
				t.Errorf("destination content = %q, want %q", got, want)
			}
		})
	}
}

func TestCopyEmpty(t *testing.T) {
	for name, copy := range copiers {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "src")
			dst := filepath.Join(dir, "dst")
			if err := os.WriteFile(src, nil, 0o600); err != nil {
				t.Fatalf("write source: %v", err)
			}

			if err := copy(src, dst); err != nil {
				t.Fatalf("copy: %v", err)
			}

			fi, err := os.Stat(dst)
			if err != nil {
				t.Fatalf("stat destination: %v", err)
			}
			if fi.Size() != 0 {
				t.Errorf("destination size = %d, want 0", fi.Size())
			}
		})
	}
}

func TestMissingSourceLeavesNoDestination(t *testing.T) {
	for name, copy := range copiers {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "absent")
			dst := filepath.Join(dir, "dst")

			if err := copy(src, dst); err == nil {
				t.Fatal("copy of a missing source returned nil error")
			}

			if _, err := os.Stat(dst); !os.IsNotExist(err) {
				t.Errorf("destination exists after a failed copy, stat error = %v", err)
			}
		})
	}
}

func TestUncreatableDestination(t *testing.T) {
	for name, copy := range copiers {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "src")
			if err := os.WriteFile(src, []byte("payload"), 0o600); err != nil {
				t.Fatalf("write source: %v", err)
			}

			// A destination below a regular file cannot be created.
			dst := filepath.Join(src, "dst")
			if err := copy(src, dst); err == nil {
				t.Fatal("copy to an uncreatable destination returned nil error")
			}
		})
	}
}
