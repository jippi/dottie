package test_helpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/jippi/dottie/cmd"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RunFilebasedCommandTests(t *testing.T, globalArgs ...string) {
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
		commandArgs  []string
	}

	tests := []testData{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".command.txt") {
			continue
		}

		base := strings.TrimSuffix(file.Name(), ".command.txt")

		content, err := os.ReadFile("tests/" + file.Name())
		require.NoErrorf(t, err, "failed to read file: %s", "tests/"+file.Name())

		var commandArgs []string

		str := string(bytes.TrimFunc(content, unicode.IsSpace))
		if len(str) > 0 {
			commandArgs = strings.Split(str, "\n")
		}

		test := testData{
			name:         base,
			goldenStdout: "stdout",
			goldenStderr: "stderr",
			goldenEnv:    "env",
			envFile:      base + ".env",
			commandArgs:  commandArgs,
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

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()

			// Copy the input.env to temporary place
			err := copyFile(t, "tests/"+tt.envFile, tmpDir+"/tmp.env")
			require.NoErrorf(t, err, "failed to copy [%s] to TempDir", tt.envFile)

			// Point args to the copied temp env file
			args := []string{"-f", tmpDir + "/tmp.env"}
			args = append(args, globalArgs...)

			if len(tt.commandArgs) > 0 {
				args = append(args, tt.commandArgs...)
			}

			// Prepare output buffers
			stdout := bytes.Buffer{}
			stderr := bytes.Buffer{}

			// Prepare command
			root := cmd.NewCommand()
			root.SetArgs(args)
			root.SetOut(&stdout)
			root.SetErr(&stderr)

			// Run command
			out, err := root.ExecuteC()
			if err != nil {
				// Append errors to stderr
				stderr.WriteString(fmt.Sprintf("%+v", err))
			}

			// Assert we got a Cobra command back
			require.NotNil(t, out, "expected a return value")

			// Read the modified .env file back
			modifiedEnv, err := os.ReadFile(tmpDir + "/tmp.env")
			require.NoErrorf(t, err, "failed to read file: %s/tmp.env", tmpDir)

			// Assert stdout + stderr + modified env file is as expected
			if stdout.Len() == 0 {
				assert.NoFileExists(t, "tests/"+tt.name+"/stdout.golden")
			} else {
				golden.Assert(t, tt.goldenStdout, stdout.Bytes())
			}

			if stderr.Len() == 0 {
				assert.NoFileExists(t, "tests/"+tt.name+"/stderr.golden")
			} else {
				golden.Assert(t, tt.goldenStderr, stderr.Bytes())
			}

			golden.Assert(t, tt.goldenEnv, modifiedEnv)
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
