package ast

import (
	"fmt"
)

func ContextualError(stmt Statement, err error) error {
	if err == nil {
		return nil
	}

	if stmt == nil {
		return err
	}

	var pos *Position
	switch val := stmt.(type) {
	case *Assignment:
		pos = &val.Position
	}

	if pos == nil {
		return err
	}

	return fmt.Errorf("%w (%s)", err, pos)
}
