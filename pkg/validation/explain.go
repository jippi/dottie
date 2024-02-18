package validation

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

type multiError interface {
	Errors() []error
}

func Explain(ctx context.Context, doc *ast.Document, inputError any, keyErr *ast.ValidationError, applyFixer, showField bool) string {
	var buff bytes.Buffer

	writer := tui.NewWriter(ctx, &buff)

	dark := writer.Dark()
	bold := writer.Warning().Copy(tui.WithEmphasis(true))
	danger := writer.Danger()
	light := writer.Light()
	primary := writer.Primary()

	stderr := tui.WriterFromContext(ctx, tui.Stderr)

	switch err := inputError.(type) {
	// Unwrap the ValidationError
	case *ast.ValidationError:
		return Explain(ctx, doc, err.WrappedError, err, applyFixer, showField)

	case multiError:
		for _, e := range err.Errors() {
			buff.WriteString(Explain(ctx, doc, e, nil, applyFixer, showField))
			buff.WriteString("\n")
		}

	// user configuration error
	case validator.InvalidValidationError:
		danger.Println("invalid validation rules: " + err.Error())

	// actual validation error
	case validator.ValidationErrors:
		if showField {
			danger.Print(keyErr.Assignment.Name)

			light.Print(" (", keyErr.Assignment.Position, ")")

			dark.Println()
		}

		for _, rule := range err {
			askToFix := applyFixer

			if showField {
				primary.Print("  * ")
			}

			switch rule.ActualTag() {
			case "dir":
				light.Println("(dir) The directory [" + bold.Sprintf(keyErr.Assignment.Interpolated) + "] does not exist.")

				if askToFix {
					fmt.Fprintln(os.Stderr, buff.String())
					buff.Reset()

					AskToCreateDirectory(ctx, keyErr.Assignment.Interpolated)

					askToFix = false
				}

			case "file":
				light.Print("(file) The file [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Println("] does not exist.")

			case "oneof":
				light.Print("(oneof) The value [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Print("] is not one of [")
				bold.Print(rule.Param())
				light.Println("].")

			case "number":
				light.Print("(number) The value [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Println("] is not a valid number.")

			case "email":
				light.Print("(email) The value [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Println("] is not a valid e-mail.")

			case "required":
				light.Println("(required) This value must not be empty/blank.")

			case "fqdn":
				light.Print("(fqdn) The value [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Println("] is not a valid Fully Qualified Domain Name (FQDN).")

			case "hostname":
				light.Print("(hostname) The value [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Println("] is not a valid hostname (e.g., 'example.com').")

			case "ne":
				light.Print("(ne) The value [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Print("] must NOT be equal to [")
				bold.Print(rule.Param())
				light.Println("], please change it.")

			case "boolean":
				light.Print("(boolean) The value [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Print("] is not a valid boolean.")

			case "http_url":
				light.Print("(http_url) The value [")
				bold.Print(keyErr.Assignment.Interpolated)
				light.Println("] is not a valid HTTP URL (e.g., 'https://example.com').")

			default:
				light.Printf("(%s) The value [", rule.ActualTag())
				bold.Print(keyErr.Assignment.Interpolated)
				light.Println("] failed validation.")
			}

			if askToFix {
				stderr.NoColor().Println(buff.String())
				buff.Reset()

				AskToSetValue(ctx, doc, keyErr.Assignment)
			}
		}

	default:
		danger.Printfln("(error %T) %+s", err, err)
	}

	return buff.String()
}

func AskToCreateDirectory(ctx context.Context, path string) {
	var (
		confirm = true
		stderr  = tui.WriterFromContext(ctx, tui.Stderr)
	)

	err := huh.NewConfirm().
		Title("\nDo you want this program to create the directory for you?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm).
		Run()
	if err != nil {
		stderr.Warning().Println("    Prompt cancelled: " + err.Error())

		return
	}

	if !confirm {
		stderr.Warning().Println("    Prompt cancelled")

		return
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		stderr.Danger().Println("    Could not create directory: " + err.Error())

		return
	}

	stderr.Success().Println("    Directory was successfully created")
}

func AskToSetValue(ctx context.Context, doc *ast.Document, assignment *ast.Assignment) {
	var (
		value  string
		stderr = tui.WriterFromContext(ctx, tui.Stderr)
	)

	err := huh.NewInput().
		Title("Please provide value for " + assignment.Name).
		Description(strings.TrimSpace(assignment.Documentation(true)) + ". (Press Ctrl+C to exit/cancel)").
		Validate(func(s string) error {
			err := validator.New().Var(s, assignment.ValidationRules())
			if err != nil {
				z := &ast.ValidationError{
					WrappedError: err,
					Assignment:   assignment,
				}

				return errors.New(Explain(ctx, doc, z, z, false, false))
			}

			return nil
		}).
		Value(&value).
		Run()
	if err != nil {
		stderr.Warning().Println("    Prompt cancelled: " + err.Error())

		return
	}

	assignment.Literal = value
	if err := pkg.Save(ctx, assignment.Position.File, doc); err != nil {
		stderr.Danger().Println("    Could not update key with value [" + value + "]: " + err.Error())

		return
	}

	stderr.Success().Println("    Successfully updated key with value [" + value + "]")
}
