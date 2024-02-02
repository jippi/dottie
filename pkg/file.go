package pkg

import (
	"io"
	"os"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/parser"
	"github.com/jippi/dottie/pkg/scanner"
)

func Load(filename string) (rows *ast.Document, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()

	return Parse(r)
}

func Save(filename string, env *ast.Document) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(env.RenderFull())
	return err
}

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader) (*ast.Document, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return parser.New(scanner.New(string(input))).Parse()
}
