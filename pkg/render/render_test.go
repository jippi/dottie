package render_test

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/render"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestFormatter(t *testing.T) {
	t.Parallel()

	golden := goldie.New(
		t,
		goldie.WithFixtureDir("test-fixtures/formatter/output"),
		goldie.WithNameSuffix(".golden.env"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)

	// Build test data set
	type testData struct {
		name     string
		filename string
	}

	tests := []testData{}

	files, err := os.ReadDir("test-fixtures/formatter")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		switch {
		case strings.HasSuffix(file.Name(), ".env"):
			testName := strings.TrimSuffix(file.Name(), ".env")

			test := testData{
				name:     testName,
				filename: "test-fixtures/formatter/" + file.Name(),
			}
			tests = append(tests, test)

		default:
			require.FailNow(t, "unexpected test file: ["+file.Name()+"]")
		}
	}

	// Run tests

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env, err := pkg.Load(tt.filename)
			require.NoError(t, err)

			golden.Assert(t, tt.name, []byte(render.NewFormatter().Statement(env).String()))
		})
	}
}
