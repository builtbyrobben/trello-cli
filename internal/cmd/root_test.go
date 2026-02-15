package cmd

import (
	"testing"
)

func TestExecute_Help(t *testing.T) {
	t.Parallel()

	err := Execute([]string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_Version(t *testing.T) {
	t.Parallel()

	err := Execute([]string{"version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_InvalidCommand(t *testing.T) {
	t.Parallel()

	err := Execute([]string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}
