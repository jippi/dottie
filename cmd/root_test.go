package cmd_test

import (
	"bytes"
	"testing"

	"github.com/jippi/dottie/cmd"
	"github.com/jippi/dottie/pkg/test_helpers"
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

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)

	executed, err := cmd.RunCommand(ctx, []string{"--version"}, &stdout, &stderr)
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

func TestRunCommandUnknownSubcommand(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)

	executed, err := cmd.RunCommand(ctx, []string{"does-not-exist"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected unknown subcommand to return an error")
	}

	if executed == nil {
		t.Fatal("expected RunCommand to return executed command")
	}

	if executed.Name() != "dottie" {
		t.Fatalf("expected root command for unknown subcommand, got %q", executed.Name())
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout for unknown subcommand, got %q", stdout.String())
	}

	if !bytes.Contains(stderr.Bytes(), []byte("unknown command \"does-not-exist\" for \"dottie\"")) {
		t.Fatalf("expected unknown-command message in stderr, got %q", stderr.String())
	}

	if !bytes.Contains(stderr.Bytes(), []byte("Run 'dottie --help' for usage.")) {
		t.Fatalf("expected help hint in stderr, got %q", stderr.String())
	}
}

func TestRunCommandWithoutArgsPrintsHelp(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)

	executed, err := cmd.RunCommand(ctx, []string{}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("expected no error when running root without args, got %v", err)
	}

	if executed == nil {
		t.Fatal("expected RunCommand to return executed command")
	}

	if executed.Name() != "dottie" {
		t.Fatalf("expected root command when running without args, got %q", executed.Name())
	}

	if stdout.Len() == 0 {
		t.Fatal("expected help output on stdout when running root without args")
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr when running root without args, got %q", stderr.String())
	}
}
