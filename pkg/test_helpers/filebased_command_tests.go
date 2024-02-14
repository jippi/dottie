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

	golden := goldie.New(
		t,
		goldie.WithFixtureDir("tests"),
		goldie.WithNameSuffix(""),
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)

	files, err := os.ReadDir("tests")
	if err != nil {
		log.Fatal(err)
	}

	// Build test data set
	type testData struct {
		name         string
		directory    string
		envFile      string
		goldenStdout string
		goldenStderr string
		goldenEnv    string
		command      []string
	}

	tests := []testData{}

	for _, file := range files {
		if !file.IsDir() {
			require.FailNowf(t, "Unexpected file [%s]. Please make a sub-directory for you test", file.Name())
		}

		content, err := os.ReadFile("tests/" + file.Name() + "/input.command.txt")
		require.NoErrorf(t, err, "failed to read file: %s", "tests/"+file.Name()+"/command.txt")

		test := testData{
			name:         file.Name(),
			directory:    "tests/" + file.Name(),
			goldenStdout: file.Name() + "/golden.stdout",
			goldenStderr: file.Name() + "/golden.stderr",
			goldenEnv:    file.Name() + "/golden.env",
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
			err := Copy(tt.directory+"/input.env", tmpDir+"/tmp.env")
			require.NoError(t, err, "failed to copy [input.env] to TempDir")

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
			require.NoErrorf(t, err, "failed to read file: %s", tt.directory+"/tmp.env")

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
