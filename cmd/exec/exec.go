package exec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/template"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// RunOptions holds the parameters for the exec logic, decoupled from cobra flags.
type RunOptions struct {
	Filename         string
	ExcludeKeyPrefix []string
	IgnoreRules      []string
	Validate         bool
	Save             bool
	Verbose          bool
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exec",
		Short:   "Update the .env file from a source",
		GroupID: "manipulate",
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}

	cmd.Flags().String("source", "", "URL or local file path to the upstream source file. This will take precedence over any [@dottie/source] annotation in the file")
	cmd.Flags().StringSlice("ignore-rule", []string{}, "Ignore this validation rule (e.g. 'dir')")
	cmd.Flags().StringSlice("exclude-key-prefix", []string{}, "Ignore these KEY prefixes")

	shared.BoolWithInverse(cmd, "error-on-missing-key", false, "Error if a KEY in FILE is missing from SOURCE", "Add KEY to FILE if missing from SOURCE")
	shared.BoolWithInverse(cmd, "validate", true, "Validation errors will abort the update", "Validation errors will be printed but will not fail the update")
	shared.BoolWithInverse(cmd, "save", true, "Save the document after processing", "Do not save the document after processing")

	cmd.Flags().BoolP("verbose", "v", false, "Show detailed output including command results (default: off to avoid printing secrets)")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	verbose := shared.BoolFlag(cmd.Flags(), "verbose")
	if !cmd.Flags().Changed("verbose") {
		if env := os.Getenv("DOTTIE_VERBOSE"); env == "1" || env == "true" {
			verbose = true
		}

		if os.Getenv("DOTTIE_DEBUG") == "1" {
			verbose = true
		}
	}

	return Run(cmd.Context(), RunOptions{
		Filename:         cmd.Flag("file").Value.String(),
		ExcludeKeyPrefix: shared.StringSliceFlag(cmd.Flags(), "exclude-key-prefix"),
		IgnoreRules:      shared.StringSliceFlag(cmd.Flags(), "ignore-rule"),
		Validate:         shared.BoolWithInverseValue(cmd.Flags(), "validate"),
		Save:             shared.BoolWithInverseValue(cmd.Flags(), "save"),
		Verbose:          verbose,
	})
}

// Run executes all dottie/exec annotations on assignments in the named file.
func Run(ctx context.Context, opts RunOptions) error {
	document, err := pkg.Load(ctx, opts.Filename)
	if err != nil {
		return err
	}

	var selectors []ast.Selector

	selectors = append(selectors, ast.ExcludeDisabledAssignments)

	for _, prefix := range opts.ExcludeKeyPrefix {
		selectors = append(selectors, ast.ExcludeKeyPrefix(prefix))
	}

	out := tui.StdoutFromContext(ctx)
	errOut := tui.StderrFromContext(ctx)
	count := 0

	for _, assignment := range document.AllAssignments(selectors...) {
		annotations := assignment.Annotation("dottie/exec")
		if len(annotations) == 0 {
			continue
		}

		if len(annotations) > 1 {
			return fmt.Errorf("multiple exec annotations found for assignment [ %s ]", assignment.Name)
		}

		if opts.Verbose && count > 0 {
			out.NoColor().Println()
		}

		count++

		out.Info().Printfln("Running exec command for assignment [ %s ]", assignment.Name)

		if opts.Verbose {
			out.Dark().Printfln("  Command: [ %s ]", annotations[0])
		}

		var buf bytes.Buffer

		runner, err := interp.New(interp.StdIO(os.Stdin, &buf, errOut.GetWriter()))
		if err != nil {
			return err
		}

		runner.Env = template.EnvironmentHelper{
			Resolver:            document.InterpolationMapper(assignment),
			AccessibleVariables: document.AccessibleVariables(assignment),
			MissingKeyCallback:  template.DefaultMissingKeyCallback(ctx, assignment.Literal),
		}

		pwd, _ := os.Getwd()
		runner.Dir = pwd

		prog, err := syntax.NewParser().Parse(strings.NewReader(annotations[0]), "")
		if err != nil {
			return err
		}

		runner.Reset()

		if err := runner.Run(ctx, prog); err != nil {
			return err
		}

		// Trim the output to remove any leading and trailing newlines
		output := strings.TrimSpace(buf.String())

		if opts.Verbose {
			out.Success().Printfln("  Output : [ %s ]", output)
		}

		// Update literal
		assignment.SetLiteral(ctx, output)

		// Validate the assignment
		validationErrors, err := document.ValidateSingleAssignment(ctx, assignment, nil, opts.IgnoreRules)
		if err != nil {
			return err
		}

		if validationErrors != nil {
			fmt.Fprintln(errOut.GetWriter(), validation.Explain(ctx, document, validationErrors, assignment, false, true))

			if opts.Validate {
				return errors.New("validation failed")
			}

			out.Warning().Println("  Validation failed, but continuing because [--no-validate] was provided")

			continue
		}

		if opts.Verbose {
			out.Success().Println("  Validation succeeded")
		}
	}

	out.NoColor().Println()
	out.Success().Println("All exec commands completed successfully")

	if !opts.Save {
		out.Warning().Println("[--no-save] was provided, not saving file")

		return nil
	}

	if err := pkg.Save(ctx, opts.Filename, document); err != nil {
		return err
	}

	out.Success().Println("File successfully saved")

	return nil
}
