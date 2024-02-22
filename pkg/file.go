package pkg

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/parser"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/scanner"
)

func Load(ctx context.Context, filename string) (*ast.Document, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(ctx, file, filename)
}

func Save(ctx context.Context, filename string, doc *ast.Document) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	res := render.NewFormatter().Statement(ctx, doc)
	if res.IsEmpty() {
		return errors.New("The rendered .env file is unexpectedly 0 bytes long - please report this as a bug (unless your file is empty)")
	}

	_, err = file.WriteString(res.String())

	return err
}

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(ctx context.Context, r io.Reader, filename string) (*ast.Document, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return parser.
		New(
			scanner.New(string(input)),
			filename,
		).
		Parse(ctx)
}
