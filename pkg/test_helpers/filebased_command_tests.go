package test_helpers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/google/shlex"
	"github.com/jippi/dottie/cmd"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Setting is a bitmask for controlling Upsert behavior
type Setting int

const (
	ReadOnly Setting = 1 << iota
)

// Has checks if [check] exists in the [settings] bitmask or not.
func (bitmask Setting) Has(setting Setting) bool {
	// If [settings] is 0, its an initialized/unconfigured bitmask, so no settings exists.
	//
	// This is true since all UpsertSetting starts from "1", not "0".
	if bitmask == 0 {
		return false
	}

	return bitmask&setting != 0
}

func RunFileBasedCommandTests(t *testing.T, settings Setting, globalArgs ...string) {
	t.Helper()

	files, err := os.ReadDir("tests")
	if err != nil {
		require.NoError(t, err, "could not read the tests/ directory")
	}

	// Build test data set
	type testData struct {
		name         string
		envFile      string
		goldenStdout string
		goldenStderr string
		goldenEnv    string
		commandsFile string
		commands     [][]string
	}

	tests := []testData{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".run") {
			continue
		}

		base := strings.TrimSuffix(file.Name(), ".run")

		content, err := os.ReadFile("tests/" + file.Name())
		require.NoErrorf(t, err, "failed to read file: %s", "tests/"+file.Name())

		var commands [][]string

		str := string(bytes.TrimFunc(content, unicode.IsSpace))
		if len(str) > 0 {
			commandArgs := strings.Split(str, "\n")
			for _, commandStr := range commandArgs {
				command, err := shlex.Split(commandStr)

				require.NoError(t, err)

				commands = append(commands, command)
			}
		}

		if len(commands) == 0 {
			commands = append(commands, []string{})
		}

		test := testData{
			name:         base,
			goldenStdout: "stdout",
			goldenStderr: "stderr",
			goldenEnv:    "env",
			envFile:      base + ".env",
			commands:     commands,
			commandsFile: "tests/" + file.Name(),
		}

		tests = append(tests, test)
	}

	// Run tests

	golden := goldie.New(
		t,
		goldie.WithFixtureDir("tests/"),
		goldie.WithSubTestNameForDir(true),
		goldie.WithNameSuffix(".golden"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)

	const sep = "-"

	for _, tt := range tests {
		tt := tt

		header := strings.Repeat(sep, 80) + "\n"
		footer := "\n" + strings.Repeat(sep, 80) + "\n\n"

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dotEnvFile := "tests/" + tt.envFile

			if !settings.Has(ReadOnly) {
				dotEnvFile = t.TempDir() + "/tmp.env"

				if _, err := os.Stat("tests/" + tt.envFile); errors.Is(err, os.ErrNotExist) {
					// Create a temporary empty .env file
					_, err := os.Create(dotEnvFile)
					require.NoErrorf(t, err, "failed to create empty .env file [ %s ] in TempDir", tt.envFile)
				} else {
					// Copy the input.env to temporary place
					err := copyFile(t, "tests/"+tt.envFile, dotEnvFile)
					require.NoErrorf(t, err, "failed to copy [ %s ] to TempDir", tt.envFile)
				}
			}

			// Prepare output buffers
			combinedStdout := bytes.Buffer{}
			combinedStderr := bytes.Buffer{}

			for idx, command := range tt.commands {
				// Point args to the copied temp env file
				args := []string{}
				args = append(args, globalArgs...)
				args = append(args, command...)

				t.Logf("Running step from line %d: %+v", idx+1, args)

				combinedStdout.WriteString(header)
				combinedStdout.WriteString(fmt.Sprintf("%s Output of command from line %d in [%s]:\n%s %+v", sep, idx+1, tt.commandsFile, sep, args))
				combinedStdout.WriteString(footer)

				combinedStderr.WriteString(header)
				combinedStderr.WriteString(fmt.Sprintf("%s Output of command from line %d in [%s]:\n%s %+v", sep, idx+1, tt.commandsFile, sep, args))
				combinedStderr.WriteString(footer)

				commandArgs := append(args, "--file", dotEnvFile)

				// Run command
				var (
					stdout = bytes.Buffer{}
					stderr = bytes.Buffer{}
					ctx    = CreateContext(t, &stdout, &stderr)
				)

				out, _ := cmd.RunCommand(ctx, commandArgs, &stdout, &stderr)

				if stdout.Len() == 0 {
					stdout.WriteString("(no output to stdout)\n")
				}

				if stderr.Len() == 0 {
					stderr.WriteString("(no output to stderr)\n")
				}

				stdout.WriteTo(&combinedStdout)
				stderr.WriteTo(&combinedStderr)

				if idx == 0 {
					header = "\n" + header
				}

				// Assert we got a Cobra command back
				require.NotNil(t, out, "expected a return value")
			}

			// Assert stdout + stderr + modified env file is as expected
			golden.Assert(t, tt.goldenStdout, combinedStdout.Bytes())
			golden.Assert(t, tt.goldenStderr, combinedStderr.Bytes())

			if !settings.Has(ReadOnly) {
				// Read the modified .env file back
				modifiedEnv, err := os.ReadFile(dotEnvFile)

				require.NoErrorf(t, err, "failed to read file: %s", dotEnvFile)
				golden.Assert(t, tt.goldenEnv, modifiedEnv)
			} else {
				assert.NoFileExists(t, "tests/"+tt.name+"/env.golden")
			}
		})
	}
}

func copyFile(t *testing.T, src, dst string) error {
	t.Helper()

	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	info, err := srcF.Stat()
	if err != nil {
		return err
	}

	dstF, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}
