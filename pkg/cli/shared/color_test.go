package shared

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestColorEnabled(t *testing.T) {
	t.Run("default true", func(t *testing.T) {
		flags := testColorFlags(t)
		require.True(t, ColorEnabled(flags, "color"))
	})

	t.Run("explicit no-color", func(t *testing.T) {
		flags := testColorFlags(t)
		require.NoError(t, flags.Set("no-color", "true"))
		require.False(t, ColorEnabled(flags, "color"))
	})

	t.Run("explicit color false", func(t *testing.T) {
		flags := testColorFlags(t)
		require.NoError(t, flags.Set("color", "false"))
		require.False(t, ColorEnabled(flags, "color"))
	})

	t.Run("no color env wins", func(t *testing.T) {
		t.Setenv("NO_COLOR", "1")

		flags := testColorFlags(t)
		require.NoError(t, flags.Set("color", "true"))
		require.False(t, ColorEnabled(flags, "color"))
	})
}

func testColorFlags(t *testing.T) *pflag.FlagSet {
	t.Helper()

	cmd := &cobra.Command{Use: "test"}
	BoolWithInverse(cmd, "color", true, "", "")

	return cmd.Flags()
}
