package main_test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanningErrorsBuild(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "golox")

	build := exec.Command("go", "build", "-o", bin, ".")
	buildOutput, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, buildOutput)
	}

	cmd := exec.Command(bin, "./testscripts/4_scanning_errors.lox")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected command to fail")
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T", err)
	}

	if exitErr.ExitCode() != 65 {
		t.Fatalf("expected exit code 65, got %d\noutput:\n%s", exitErr.ExitCode(), output)
	}

	got := string(output)

	if !strings.Contains(got, "unexpected character: '&'") {
		t.Errorf("missing expected error message")
	}

	if !strings.Contains(got, "unexpected character: '|'") {
		t.Errorf("missing expected error message")
	}

	if !strings.Contains(got, "unterminated string") {
		t.Errorf("missing expected error message")
	}
}

func TestCLI(t *testing.T) {
	var tests = []struct {
		file     string
		expected string
	}{
		{
			file:     "./testscripts/8_vars.lox",
			expected: "42",
		},
		{
			file:     "./testscripts/9_controlflow_ifelse.lox",
			expected: "yes, true, indeed.",
		},
		{
			file:     "./testscripts/8_blocks.lox",
			expected: "inner a\nouter b\nglobal c\nouter a\nouter b\nglobal c\nglobal a\nglobal b\nglobal c",
		},
		{
			file:     "./testscripts/9_controlflow_andor.lox",
			expected: "first yes\nsecond yes\nthird yes",
		},
	}

	var failed bool

	for _, test := range tests {

		cmd := exec.Command(
			"go",
			"run",
			".",
			test.file,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("executing script %q failed: %v", test.file, err)
			failed = true
			continue
		}

		got := strings.TrimSpace(string(output))
		want := test.expected

		if got != want {
			t.Logf("failed %s: got %q, want %q", test.file, got, want)
			failed = true
		}
	}

	if failed {
		t.Fatal("failed TestCLI")
	}
}
