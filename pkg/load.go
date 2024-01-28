package pkg

import (
	"fmt"
	"io"
	"os"

	"dotfedi/pkg/ast"
	"dotfedi/pkg/parser"
	"dotfedi/pkg/scanner"
)

func Load(filename string) (rows *ast.File, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	return Parse(f)
}

func Save(filename string, env *ast.File) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(env.Render())
	return err
}

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader) (*ast.File, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	s := scanner.New(string(input))
	p := parser.New(s)

	statement, err := p.Parse()
	if err != nil {
		return nil, err
	}

	fileStmt, ok := statement.(*ast.File)
	if !ok {
		return nil, fmt.Errorf("(A) unexpected statement: %T", statement)
	}

	return fileStmt, nil
}
