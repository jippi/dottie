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

	g := goldie.New(
		t,
		goldie.WithFixtureDir("test-fixtures/formatter"),
		goldie.WithNameSuffix(".golden.env"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)

	files, err := os.ReadDir("test-fixtures/formatter")
	if err != nil {
		log.Fatal(err)
	}

	// Build test data set
	type testData struct {
		name     string
		filename string
	}

	tests := []testData{}

	for _, file := range files {
		switch {
		case strings.HasSuffix(file.Name(), ".input.env"):
			testName := strings.TrimSuffix(file.Name(), ".input.env")

			tests = append(tests, testData{name: testName, filename: "test-fixtures/formatter/" + file.Name()})

		case strings.HasSuffix(file.Name(), ".golden.env"):
		default:
			panic("unexpected file")
		}
	}

	// Run tests

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env, err := pkg.Load(tt.filename)
			require.NoError(t, err)

			g.Assert(t, tt.name, []byte(render.NewFormatter().Statement(env)))
		})
	}
}
