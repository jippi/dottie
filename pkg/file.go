package pkg

import (
	"fmt"
	"io"
	"os"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/parser"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/scanner"
)

func Load(filename string) (doc *ast.Document, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()

	return Parse(r, filename)
}

func Save(filename string, doc *ast.Document) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	res := render.NewFormatter().Statement(doc)
	if res.Empty() {
		return fmt.Errorf("The rendered .env file is unexpectedly 0 bytes long - please report this as a bug (unless your file is empty)")
	}

	_, err = f.WriteString(res.GetWithEOF())

	return err
}

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader, filename string) (*ast.Document, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return parser.
		New(
			scanner.New(string(input)),
			filename,
		).
		Parse()
}
