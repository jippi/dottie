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
		goldie.WithTestNameForDir(false),
	)

	files, err := os.ReadDir("test-fixtures/formatter")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		switch {
		case strings.HasSuffix(file.Name(), ".input.env"):
			env, err := pkg.Load("test-fixtures/formatter/" + file.Name())
			require.NoError(t, err)

			testName := strings.TrimSuffix(file.Name(), ".input.env")

			g.Assert(t, testName, []byte(render.NewFormatter(env)))

		case strings.HasSuffix(file.Name(), ".golden.env"):
		default:
			panic("unexpected file")
		}
	}
}
