// Copyright 2026 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestServerAcceptsEEPQ runs the built server binary and checks that it
// lists the -eepq flag. On 2026-11-15 a server built from upspin before
// eepq stopped on that flag in production while every test passed.
func TestServerAcceptsEEPQ(t *testing.T) {
	bin := os.Getenv("UPSPINSERVER_DRIVE")
	if bin == "" {
		t.Fatal("UPSPINSERVER_DRIVE is not set; run this test with bazel test")
	}
	// The path is relative to the runfiles of the main repository.
	if dir := os.Getenv("TEST_SRCDIR"); dir != "" {
		bin = filepath.Join(dir, "_main", bin)
	}
	// -help prints the flags and exits with status 2; only a failure to
	// start the program is an error here.
	out, err := exec.Command(bin, "-help").CombinedOutput()
	if _, ok := err.(*exec.ExitError); err != nil && !ok {
		t.Fatalf("run %s: %v", bin, err)
	}
	if !strings.Contains(string(out), "-eepq") {
		t.Fatalf("%s -help does not list -eepq:\n%s", bin, out)
	}
}
