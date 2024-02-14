package test_helpers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/jippi/dottie/cmd"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func RunFilebasedCommandTests(t *testing.T) {
	t.Helper()

	files, err := os.ReadDir("tests")
	if err != nil {
		log.Fatal(err)
	}

	// Build test data set
	type testData struct {
		name         string
		envFile      string
		goldenStdout string
		goldenStderr string
		goldenEnv    string
		command      []string
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

		test := testData{
			name:         base,
			goldenStdout: "stdout",
			goldenStderr: "stderr",
			goldenEnv:    "env",
			envFile:      file.Name(),
			command:      strings.Split(strings.TrimSpace(string(content)), "\n"),
		}

		tests = append(tests, test)
	}

	// Run tests

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()

			// Copy the input.env to temporary place
			err := Copy("tests/"+tt.envFile, tmpDir+"/tmp.env")
			require.NoErrorf(t, err, "failed to copy [%s] to TempDir", tt.envFile)

			// Point args to the copied temp env file
			args := []string{"-f", tmpDir + "/tmp.env"}
			args = append(args, tt.command...)

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

			golden := goldie.New(
				t,
				goldie.WithFixtureDir("tests/"),
				// goldie.WithTestNameForDir(false),
				goldie.WithSubTestNameForDir(true),
				goldie.WithNameSuffix(".golden"),
				goldie.WithDiffEngine(goldie.ColoredDiff),
			)

			// Assert stdout + stderr + modified env file is as expected
			golden.Assert(t, tt.goldenStdout, stdout.Bytes())
			golden.Assert(t, tt.goldenStderr, stderr.Bytes())
			golden.Assert(t, tt.goldenEnv, modifiedEnv)
		})
	}
}

func Copy(src, dst string) error {
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
