package validation

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

func Explain(env *ast.Document, keyErr ValidationError) {
	normal := tui.Theme.Dark.StderrPrinter()
	bold := tui.Theme.Dark.StderrPrinter(tui.WithEmphasis(true))
	danger := tui.Theme.Danger.StderrPrinter()
	light := tui.Theme.Light.StderrPrinter()
	secondary := tui.Theme.Primary.StderrPrinter()

	switch err := keyErr.Error.(type) {
	// user configuration error
	case validator.InvalidValidationError:
		danger.Println("invalid validation rules: " + err.Error())

	// actual validation error
	case validator.ValidationErrors:
		danger.Print(keyErr.Assignment.Name)

		light.Print(" (", keyErr.Assignment.Position, ")")

		normal.Println()

		for _, rule := range err {
			secondary.Print("  * ")

			switch rule.ActualTag() {
			case "dir":
				light.Println("(dir) The directory [" + bold.Sprintf(keyErr.Assignment.Interpolated) + "] does not exist.")
				AskToCreateDirectory(keyErr.Assignment.Interpolated)

			case "file":
				light.Println("(file) The file [" + bold.Sprintf(keyErr.Assignment.Interpolated) + "] does not exist.")

			case "oneof":
				light.Println("(oneof) The value [" + bold.Sprintf(keyErr.Assignment.Interpolated) + "] must be one of [" + rule.Param() + "]")

			case "email":
				light.Println("(email) Expected a valid e-mail, but got [" + bold.Sprintf(keyErr.Assignment.Interpolated) + "].")
				AskToSetValue(env, keyErr.Assignment)

			case "required":
				light.Println("(required) This key must not have an empty value.")

			case "fqdn":
				light.Println("(fqdn) This key must have a valid hostname.")
				AskToSetValue(env, keyErr.Assignment)

			default:
				light.Printfln("(%s) failed validation", rule.ActualTag())
			}
		}

		normal.Println()

	case error:
		light.Printfln("(error) %s", err)

	default:
		panic(fmt.Errorf("unknown error type for field type: %T", err))
	}
}

func AskToCreateDirectory(path string) {
	confirm := true

	err := huh.NewConfirm().
		Title("\nDo you want this program to create the directory for you?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm).
		Run()
	if err != nil {
		tui.Theme.Warning.StderrPrinter().Println("    Prompt cancelled: " + err.Error())

		return
	}

	if !confirm {
		tui.Theme.Warning.StderrPrinter().Println("    Prompt cancelled")

		return
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		tui.Theme.Danger.StderrPrinter().Println("    Could not create directory: " + err.Error())

		return
	}

	tui.Theme.Success.StderrPrinter().Println("    Directory was successfully created")
}

func AskToSetValue(env *ast.Document, assignment *ast.Assignment) {
	var value string

	err := huh.NewInput().
		Title("Please provide input").
		Description(strings.TrimSpace(assignment.Documentation(true))).
		Validate(func(s string) error {
			return validator.New().Var(s, assignment.ValidationRules())
		}).
		Value(&value).
		Run()
	if err != nil {
		tui.Theme.Warning.StderrPrinter().Println("    Prompt cancelled: " + err.Error())

		return
	}

	assignment.Literal = value
	if err := pkg.Save(assignment.Position.File, env); err != nil {
		tui.Theme.Danger.StderrPrinter().Println("    Could not update key with value [" + value + "]: " + err.Error())

		return
	}

	tui.Theme.Success.StderrPrinter().Println("    Successfully updated key with value [" + value + "]")
}
