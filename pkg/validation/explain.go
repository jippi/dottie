package validation

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/huh"
	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

type multiError interface {
	Errors() []error
}

func Explain(ctx context.Context, doc *ast.Document, inputError any, assignment *ast.Assignment, applyFixer, showField bool) string {
	var buff bytes.Buffer

	writer := tui.NewWriter(ctx, &buff)

	dark := writer.NoColor()
	bold := writer.Warning().Copy(tui.WithEmphasis(true))
	danger := writer.Danger()
	light := writer.NoColor()
	primary := writer.Primary()

	stderr := tui.WriterFromContext(ctx, tui.Stderr)

	switch err := inputError.(type) {
	// Unwrap the ValidationError
	case *ast.ValidationError:
		return Explain(ctx, doc, err.WrappedError, assignment, applyFixer, showField)

		// Unwrap a list of validation errors
	case ast.ValidationErrors:
		if showField {
			danger.Print("  ", assignment.Name)
			dark.Println(" (", assignment.Position, ")")
		}

		for _, e := range err.Errors() {
			buff.WriteString(strings.TrimRightFunc(Explain(ctx, doc, e, assignment, applyFixer, false), unicode.IsSpace))
			buff.WriteString("\n")
			buff.WriteString("\n")
		}

		x := strings.TrimRightFunc(buff.String(), unicode.IsSpace)
		buff.Reset()
		buff.WriteString(x)
		buff.WriteString("\n")

	case multiError:
		for _, e := range err.Errors() {
			buff.WriteString(Explain(ctx, doc, e, assignment, applyFixer, showField))
			buff.WriteString("\n")
		}

	// user configuration error
	case validator.InvalidValidationError:
		danger.Println("invalid validation rules: " + err.Error())

	// actual validation error
	case validator.ValidationErrors:
		if showField {
			danger.Print(assignment.Name)

			dark.Print(" (", assignment.Position, ")")

			dark.Println()
		}

		for _, rule := range err {
			askToFix := applyFixer

			primary.Print("    * ")

			if rule.Field() != assignment.Name {
				dark.Print("Field ")
				danger.Print(rule.Field())
				dark.Println(" which is dependent on this KEY failed validation")

				primary.Print("      ")
			}

			tag := rule.ActualTag()
			light.Print("(", tag, ") ")

			for _, segment := range splitHighlightedMessage(explainRuleMessage(tag, rule.Param(), assignment.Interpolated)) {
				if segment.highlighted {
					bold.Print(segment.text)

					continue
				}

				light.Print(segment.text)
			}

			light.Println()

			if tag == "dir" && askToFix {
				fmt.Fprintln(tui.StderrFromContext(ctx).NoColor(), buff.String())
				buff.Reset()

				AskToCreateDirectory(ctx, assignment.Interpolated)

				askToFix = false
			}

			if askToFix {
				stderr.NoColor().Println(buff.String())
				buff.Reset()

				AskToSetValue(ctx, doc, assignment)
			}
		}

	case error:
		danger.Printfln("%+s", err)

	default:
		danger.Printfln("(error %T) %+s", err, err)
	}

	return buff.String()
}

type messageSegment struct {
	text        string
	highlighted bool
}

func splitHighlightedMessage(message string) []messageSegment {
	segments := []messageSegment{}

	rest := message

	for len(rest) > 0 {
		open := strings.Index(rest, "[")
		if open == -1 {
			segments = append(segments, messageSegment{text: rest})

			break
		}

		if open > 0 {
			segments = append(segments, messageSegment{text: rest[:open]})
		}

		remaining := rest[open+1:]

		closeIndex := strings.Index(remaining, "]")
		if closeIndex == -1 {
			segments = append(segments, messageSegment{text: rest[open:]})

			break
		}

		segments = append(segments, messageSegment{text: "["})
		segments = append(segments, messageSegment{text: remaining[:closeIndex], highlighted: true})
		segments = append(segments, messageSegment{text: "]"})

		rest = remaining[closeIndex+1:]
	}

	return segments
}

func explainRuleMessage(tag, param, value string) string {
	switch tag {
	case "required":
		return "This value is required and cannot be empty."
	case "omitempty":
		return "The value is only validated when it is non-empty, and the provided value did not pass the remaining rule(s)."
	case "required_if":
		return fmt.Sprintf("This value is required when [%s].", param)
	case "required_unless":
		return fmt.Sprintf("This value is required unless [%s].", param)
	case "required_with":
		return fmt.Sprintf("This value is required when any of [%s] is set.", param)
	case "required_with_all":
		return fmt.Sprintf("This value is required when all of [%s] are set.", param)
	case "required_without":
		return fmt.Sprintf("This value is required when any of [%s] is missing.", param)
	case "required_without_all":
		return fmt.Sprintf("This value is required when all of [%s] are missing.", param)
	case "excluded_if":
		return fmt.Sprintf("This value must be empty when [%s].", param)
	case "excluded_unless":
		return fmt.Sprintf("This value must be empty unless [%s].", param)

	case "len":
		return fmt.Sprintf("The value [%s] must have exact length/value [%s].", value, param)
	case "min":
		return fmt.Sprintf("The value [%s] must be at least [%s].", value, param)
	case "max":
		return fmt.Sprintf("The value [%s] must be at most [%s].", value, param)
	case "eq":
		return fmt.Sprintf("The value [%s] must be exactly [%s].", value, param)
	case "ne":
		return fmt.Sprintf("The value [%s] must not be equal to [%s].", value, param)
	case "gt":
		return fmt.Sprintf("The value [%s] must be greater than [%s].", value, param)
	case "gte":
		return fmt.Sprintf("The value [%s] must be greater than or equal to [%s].", value, param)
	case "lt":
		return fmt.Sprintf("The value [%s] must be less than [%s].", value, param)
	case "lte":
		return fmt.Sprintf("The value [%s] must be less than or equal to [%s].", value, param)
	case "oneof":
		return fmt.Sprintf("The value [%s] must be one of [%s].", value, param)
	case "oneofci":
		return fmt.Sprintf("The value [%s] must case-insensitively match one of [%s].", value, param)

	case "number":
		return fmt.Sprintf("The value [%s] is not a valid number.", value)
	case "numeric":
		return fmt.Sprintf("The value [%s] must be a numeric string.", value)
	case "boolean":
		return fmt.Sprintf("The value [%s] is not a valid boolean.", value)
	case "alpha":
		return fmt.Sprintf("The value [%s] must contain only letters.", value)
	case "alphanum":
		return fmt.Sprintf("The value [%s] must contain only letters and digits.", value)
	case "ascii":
		return fmt.Sprintf("The value [%s] must contain only ASCII characters.", value)
	case "lowercase":
		return fmt.Sprintf("The value [%s] must be all lowercase.", value)
	case "uppercase":
		return fmt.Sprintf("The value [%s] must be all uppercase.", value)
	case "contains":
		return fmt.Sprintf("The value [%s] must contain [%s].", value, param)
	case "excludes":
		return fmt.Sprintf("The value [%s] must not contain [%s].", value, param)
	case "startswith":
		return fmt.Sprintf("The value [%s] must start with [%s].", value, param)
	case "endswith":
		return fmt.Sprintf("The value [%s] must end with [%s].", value, param)

	case "email":
		return fmt.Sprintf("The value [%s] is not a valid e-mail address.", value)
	case "url":
		return fmt.Sprintf("The value [%s] is not a valid URL.", value)
	case "uri":
		return fmt.Sprintf("The value [%s] is not a valid URI.", value)
	case "http_url":
		return fmt.Sprintf("The value [%s] is not a valid HTTP/HTTPS URL.", value)
	case "https_url":
		return fmt.Sprintf("The value [%s] is not a valid HTTPS URL.", value)
	case "hostname":
		return fmt.Sprintf("The value [%s] is not a valid hostname.", value)
	case "hostname_rfc1123":
		return fmt.Sprintf("The value [%s] is not a valid RFC1123 hostname.", value)
	case "fqdn":
		return fmt.Sprintf("The value [%s] is not a valid fully qualified domain name (FQDN).", value)
	case "hostname_port":
		return fmt.Sprintf("The value [%s] is not a valid hostname:port pair.", value)
	case "ip":
		return fmt.Sprintf("The value [%s] is not a valid IP address.", value)
	case "ipv4":
		return fmt.Sprintf("The value [%s] is not a valid IPv4 address.", value)
	case "ipv6":
		return fmt.Sprintf("The value [%s] is not a valid IPv6 address.", value)
	case "cidr":
		return fmt.Sprintf("The value [%s] is not a valid CIDR block.", value)
	case "mac":
		return fmt.Sprintf("The value [%s] is not a valid MAC address.", value)
	case "dir":
		return fmt.Sprintf("The directory [%s] does not exist.", value)
	case "dirpath":
		return fmt.Sprintf("The value [%s] is not a valid directory path.", value)
	case "file":
		return fmt.Sprintf("The file [%s] does not exist.", value)
	case "filepath":
		return fmt.Sprintf("The value [%s] is not a valid file path.", value)

	case "uuid":
		return fmt.Sprintf("The value [%s] is not a valid UUID.", value)
	case "ulid":
		return fmt.Sprintf("The value [%s] is not a valid ULID.", value)
	case "semver":
		return fmt.Sprintf("The value [%s] is not a valid semantic version.", value)
	case "cron":
		return fmt.Sprintf("The value [%s] is not a valid cron expression.", value)
	case "json":
		return fmt.Sprintf("The value [%s] is not valid JSON.", value)
	case "jwt":
		return fmt.Sprintf("The value [%s] is not a valid JWT.", value)
	case "hexcolor":
		return fmt.Sprintf("The value [%s] is not a valid hex color.", value)
	case "rgb":
		return fmt.Sprintf("The value [%s] is not a valid RGB color.", value)
	case "rgba":
		return fmt.Sprintf("The value [%s] is not a valid RGBA color.", value)
	case "base64":
		return fmt.Sprintf("The value [%s] is not a valid base64 string.", value)
	case "timezone":
		return fmt.Sprintf("The value [%s] is not a valid time zone identifier.", value)
	default:
		return fmt.Sprintf("The value [%s] failed validation.", value)
	}
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
				z := ast.NewError(assignment, err)

				return errors.New(Explain(ctx, doc, z, assignment, false, false))
			}

			return nil
		}).
		Value(&value).
		Run()
	if err != nil {
		stderr.Warning().Println("    Prompt cancelled: " + err.Error())

		return
	}

	assignment.SetLiteral(ctx, value)

	if err := pkg.Save(ctx, assignment.Position.File, doc); err != nil {
		stderr.Danger().Println("    Could not update key with value [" + value + "]: " + err.Error())

		return
	}

	stderr.Success().Println("    Successfully updated key with value [" + value + "]")
}
