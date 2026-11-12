// Copyright 2026 The Vesper Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestRunVersionFormats(t *testing.T) {
	old := version
	version = "v1.2.3"
	t.Cleanup(func() { version = old })

	t.Run("default output is unchanged", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if got := run([]string{"version"}, &stdout, &stderr, emptyEnv); got != exitOK {
			t.Fatalf("exit code = %d, want %d", got, exitOK)
		}
		want := fmt.Sprintf("vesperctl version v1.2.3 (%s %s/%s)\n",
			runtime.Version(), runtime.GOOS, runtime.GOARCH)
		if stdout.String() != want {
			t.Errorf("stdout = %q, want %q", stdout.String(), want)
		}
	})

	t.Run("json", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if got := run([]string{"version", "-o", "json"}, &stdout, &stderr, emptyEnv); got != exitOK {
			t.Fatalf("exit code = %d, want %d (stderr: %s)", got, exitOK, stderr.String())
		}
		var got versionInfo
		if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
			t.Fatalf("stdout is not JSON: %v (%q)", err, stdout.String())
		}
		want := versionInfo{
			Version:   "v1.2.3",
			GoVersion: runtime.Version(),
			Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		}
		if got != want {
			t.Errorf("json = %+v, want %+v", got, want)
		}
		// The field names are the contract a script reads, so check them
		// literally rather than through the struct tags alone.
		for _, field := range []string{`"version"`, `"goVersion"`, `"platform"`} {
			if !strings.Contains(stdout.String(), field) {
				t.Errorf("stdout = %q, want it to contain %s", stdout.String(), field)
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if got := run([]string{"version", "--output", "yaml"}, &stdout, &stderr, emptyEnv); got != exitOK {
			t.Fatalf("exit code = %d, want %d (stderr: %s)", got, exitOK, stderr.String())
		}
		if want := "version: v1.2.3"; !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout = %q, want it to contain %q", stdout.String(), want)
		}
	})

	t.Run("unknown format", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if got := run([]string{"version", "-o", "xml"}, &stdout, &stderr, emptyEnv); got != exitUsage {
			t.Fatalf("exit code = %d, want %d", got, exitUsage)
		}
		if want := "valid formats are table, json, yaml"; !strings.Contains(stderr.String(), want) {
			t.Errorf("stderr = %q, want it to contain %q", stderr.String(), want)
		}
	})
}
