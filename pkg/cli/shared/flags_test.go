package shared_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestBoolFlag(t *testing.T) {
	t.Parallel()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.Bool("foo", true, "usage")

	if !shared.BoolFlag(flags, "foo") {
		t.Error("expected BoolFlag to return true for set flag")
	}
}

func TestStringFlag(t *testing.T) {
	t.Parallel()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("bar", "baz", "usage")

	if got := shared.StringFlag(flags, "bar"); got != "baz" {
		t.Errorf("expected StringFlag to return 'baz', got '%s'", got)
	}
}

func TestStringSliceFlag(t *testing.T) {
	t.Parallel()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.StringSlice("baz", []string{"a", "b"}, "usage")
	got := shared.StringSliceFlag(flags, "baz")

	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("expected StringSliceFlag to return [a b], got %v", got)
	}
}

func TestBoolWithInverseValue(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{}
	shared.BoolWithInverse(cmd, "foo", true, "usage", "not foo")

	flags := cmd.Flags()

	if !shared.BoolWithInverseValue(flags, "foo") {
		t.Error("expected BoolWithInverseValue to return true for set flag")
	}

	if err := flags.Set("no-foo", "true"); err != nil {
		t.Fatalf("expected setting no-foo to succeed: %v", err)
	}

	if shared.BoolWithInverseValue(flags, "foo") {
		t.Error("expected BoolWithInverseValue to return false when no-foo=true")
	}
}

func TestBoolWithInverse(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{}
	shared.BoolWithInverse(cmd, "foo", true, "usage", "not foo")
	flag := cmd.Flags().Lookup("foo")

	if flag == nil {
		t.Error("expected flag 'foo' to be registered")
	}
}
