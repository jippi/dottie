package cmd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/jippi/dottie/cmd"
)

func TestNewRootCommandRegistersCoreCommands(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCommand()

	commandNames := map[string]bool{}

	for _, sub := range root.Commands() {
		commandNames[sub.Name()] = true
	}

	for _, expected := range []string{"set", "update", "fmt", "disable", "enable", "exec", "shell", "print", "validate", "value", "groups", "json", "template"} {
		if !commandNames[expected] {
			t.Fatalf("expected root command to register %q", expected)
		}
	}
}

func TestRunCommandVersion(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	executed, err := cmd.RunCommand(context.Background(), []string{"--version"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("expected no error from --version, got %v", err)
	}

	if executed == nil {
		t.Fatal("expected RunCommand to return executed command")
	}

	if stdout.Len() == 0 {
		t.Fatal("expected --version to write to stdout")
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr for --version, got %q", stderr.String())
	}
}
