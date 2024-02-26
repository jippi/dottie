package exec

import (
	"bytes"
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

func NewCommand() *cobra.Command {
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

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	document, err := pkg.Load(cmd.Context(), filename)
	if err != nil {
		return err
	}

	out := tui.StdoutFromContext(cmd.Context())
	count := 0

	for _, assignment := range document.AllAssignments(ast.ExcludeDisabledAssignments) {
		annotations := assignment.Annotation("dottie/exec")
		if len(annotations) == 0 {
			continue
		}

		if len(annotations) > 1 {
			return fmt.Errorf("multiple exec annotations found for assignment [ %s ]", assignment.Name)
		}

		if count > 0 {
			out.NoColor().Println()
		}

		count++

		out.Info().Printfln("Running exec command for assignment [ %s ]", assignment.Name)
		out.Dark().Printfln("  Command: [ %s ]", annotations[0])

		var buf bytes.Buffer

		runner, err := interp.New(interp.StdIO(cmd.InOrStdin(), &buf, cmd.ErrOrStderr()))
		if err != nil {
			return err
		}

		runner.Env = template.EnvironmentHelper{
			Resolver:            document.InterpolationMapper(assignment),
			AccessibleVariables: document.AccessibleVariables(assignment),
			MissingKeyCallback:  template.DefaultMissingKeyCallback(cmd.Context(), assignment.Literal),
		}

		pwd, _ := os.Getwd()
		runner.Dir = pwd

		prog, err := syntax.NewParser().Parse(strings.NewReader(annotations[0]), "")
		if err != nil {
			return err
		}

		runner.Reset()

		if err := runner.Run(cmd.Context(), prog); err != nil {
			return err
		}

		// Trim the output to remove any leading and trailing newlines
		output := strings.TrimSpace(buf.String())

		out.Success().Printfln("  Output : [ %s ]", output)

		// Update literal
		assignment.SetLiteral(cmd.Context(), output)

		// Validate the assignment
		validationErrors, err := document.ValidateSingleAssignment(cmd.Context(), assignment, nil, nil)
		if err != nil {
			return err
		}

		if validationErrors != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), validation.Explain(cmd.Context(), document, validationErrors, assignment, false, true))

			return err
		}

		out.Success().Println("  Validation succeeded")
	}

	out.NoColor().Println()
	out.Success().Println("All exec commands completed successfully")

	if err := pkg.Save(cmd.Context(), filename, document); err != nil {
		return err
	}

	out.Success().Println("File successfully saved")

	return nil
}
