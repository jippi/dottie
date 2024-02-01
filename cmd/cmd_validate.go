package main

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v3"
)

var validateCommand = &cli.Command{
	Name:   "validate",
	Usage:  "Validate .env file",
	Before: setup,
	Action: func(_ context.Context, _ *cli.Command) error {
		res := env.Validate()
		if len(res) == 0 {
			fmt.Println("all god")
			return nil
		}

		for field, errIsh := range res {
			switch err := errIsh.(type) {
			// user configuration error
			case validator.InvalidValidationError:
				fmt.Println("[", field, "] invalid validation rules", err.Error())

			// actual validation error
			case validator.ValidationErrors:
				for _, rule := range err {
					fmt.Println("Field [", field, "] failed validation rule [", rule.ActualTag(), "]", rule.Param())
				}

			default:
				panic("unknown error type for field [" + field + "]")
			}
		}

		return fmt.Errorf("validation error")
	},
}
