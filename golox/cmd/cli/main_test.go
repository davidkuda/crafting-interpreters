package main_test

import (
	"os/exec"
	"strings"
	"testing"
)

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
