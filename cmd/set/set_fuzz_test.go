package set_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/davecgh/go-spew/spew"
	"github.com/jippi/dottie/cmd"
	"github.com/jippi/dottie/pkg/test_helpers"
	"github.com/jippi/dottie/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzSetCommand(f *testing.F) {
	f.SkipNow()

	f.Add("@@\v\"@23")

	f.Fuzz(doTest)
}

func TestSpecificInputs(t *testing.T) {
	t.Parallel()

	t.Run("newline", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\n")
	})

	t.Run("tab", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\t")
	})

	t.Run("slash", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\\")
	})

	t.Run("null", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\x00")
	})

	t.Run("slash-zero", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\\0")
	})

	t.Run("weird-1", func(t *testing.T) {
		t.Parallel()

		doTest(t, "@@\v\"@23")
	})

	t.Run("weird-2", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\"\n")
	})

	t.Run("weird-3", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\x00$")
	})

	t.Run("weird-4", func(t *testing.T) {
		t.Parallel()

		doTest(t, "`0|$`")
	})
}

func doTest(t *testing.T, expected string) { //nolint thelper
	// NULL bytes are acting weird when reading stdout/stderr
	// They work fine in in-memory testing, so ignoring them for now :)
	if strings.Contains(expected, "\x00") {
		t.Skip()

		return
	}

	ctx := context.TODO()

	dotEnvFile := t.TempDir() + "/tmp.env"

	_, err := os.Create(dotEnvFile)
	require.NoErrorf(t, err, "failed to create empty .env file [ %s ] in TempDir", dotEnvFile)

	t.Log("-----------------------")
	t.Log("EXPECTED VALUE")
	t.Log("-----------------------")
	dump(t, ctx, expected)

	// Set the KEY/VALUE pair
	setFailed := false

	{
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
			ctx    = test_helpers.CreateTestContext(t, &stdout, &stderr)
			args   = []string{
				"--file", dotEnvFile,
				"set",
				"--",
				"my_key", expected,
			}
		)

		t.Log("----------------------------------------------")
		t.Log("[dottie set] COMMAND")
		t.Log("----------------------------------------------")
		t.Log(spew.Sdump(args))
		t.Log()

		// Run command
		_, err := cmd.RunCommand(ctx, args, &stdout, &stderr)

		if stdout.Len() == 0 {
			stdout.WriteString("(empty)")
		}

		if stderr.Len() == 0 {
			stderr.WriteString("(empty)")
		}

		t.Log("-----------------------")
		t.Log("STDOUT:")
		t.Log("-----------------------")
		t.Log(stdout.String())
		t.Log()

		t.Log("-----------------------")
		t.Log("STDERR:")
		t.Log("-----------------------")
		t.Log(stderr.String())
		t.Log()

		switch err {
		case nil:
			out := stdout.String()

			assert.Contains(t, out, "Key [ my_key ] was successfully upserted")
			assert.Contains(t, out, "File was successfully saved")

		default:
			setFailed = true

			assert.Regexp(t, "(Invalid template)", stderr.String())
		}
	}

	// Half-way checkpoint for some checks

	{
		if setFailed {
			return
		}

		out, err := os.ReadFile(dotEnvFile)
		require.NoError(t, err)

		disk := strings.TrimRight(string(out), "\n")

		t.Log("----------------------------------------------")
		t.Log("FILE ON DISK")
		t.Log("----------------------------------------------")
		dump(t, ctx, disk)
	}

	// Read back from disk
	{
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
			ctx    = test_helpers.CreateTestContext(t, &stdout, &stderr)
			args   = []string{
				"--file", dotEnvFile,
				"value",
				"--literal",
				"my_key",
			}
		)

		t.Log("----------------------------------------------")
		t.Log("[dottie value] COMMAND")
		t.Log("---------------------------------------------")
		t.Log(spew.Sdump(args))
		t.Log()

		// Run command
		_, err := cmd.RunCommand(ctx, args, &stdout, &stderr)
		require.NoError(t, err)

		t.Log("-----------------------")
		t.Log("STDOUT:")
		t.Log("-----------------------")
		t.Log(stdout.String())

		t.Log("-----------------------")
		t.Log("STDERR:")
		t.Log("-----------------------")
		t.Log(stderr.String())

		// In cases where we get "$0" and similar, the actual interpolated output will be different than the input
		if strings.Contains(stderr.String(), "Defaulting to a blank string") {
			return
		}

		actual := stdout.String()

		t.Log("-----------------------")
		t.Log("Actual")
		t.Log("-----------------------")
		dump(t, ctx, actual)

		assert.Equal(t, fmt.Sprintf("%U", []rune(expected)), fmt.Sprintf("%U", []rune(actual)))
	}
}

func Clean(str string) string {
	return strings.Map(func(runeVal rune) rune {
		if !utf8.ValidRune(runeVal) {
			return -1
		}

		if !unicode.IsGraphic(runeVal) {
			return -1
		}

		return runeVal
	}, str)
}

func dump(t *testing.T, ctx context.Context, value string) {
	t.Helper()

	for _, line := range token.DebugStringSlice(ctx, value) {
		t.Log(line)
	}
}
