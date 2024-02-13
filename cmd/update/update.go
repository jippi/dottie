package update

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update the .env file from a source",
		GroupID: "manipulate",
		RunE:    runE,
	}

	cmd.Flags().String("source", "", "URL or local file path to the upstream source file. This will take precedence over any [@dottie/source] annotation in the file")
	shared.BoolWithInverse(cmd, "error-on-missing-key", true, "Error if a KEY in FILE is missing from SOURCE", "Add KEY to FILE if missing from SOURCE")
	shared.BoolWithInverse(cmd, "validate", true, "Validation errors will abort the update", "Validation errors will be printed but will not fail the update")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	env, err := pkg.Load(filename)
	if err != nil {
		return err
	}

	dark := tui.Theme.Dark.StdoutPrinter()
	info := tui.Theme.Info.StdoutPrinter()
	danger := tui.Theme.Danger.StdoutPrinter()
	dangerEmphasis := tui.Theme.Danger.StdoutPrinter(tui.WithEmphasis(true))
	success := tui.Theme.Success.StdoutPrinter()
	primary := tui.Theme.Primary.StdoutPrinter()

	info.Box("Starting update of " + filename + " from upstream")
	info.Println()

	dark.Println("Looking for upstream source")

	source, _ := cmd.Flags().GetString("source")
	if len(source) == 0 {
		source, err = env.GetConfig("dottie/source")
		if err != nil {
			return err
		}

		success.Println("  Found source via [dottie/source] annotation in file", primary.Sprint(filename))
	} else {
		success.Println("  Found source via CLI flag")
	}

	fmt.Println()

	dark.Println("Copying source from", primary.Sprint(source))

	tmp, err := os.CreateTemp(os.TempDir(), ".dottie.source")
	if err != nil {
		return err
	}

	// Get the pwd
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting working directory: %w", err)
	}

	// Grab source file

	client := getter.Client{
		DisableSymlinks: true,
		Mode:            getter.ClientModeFile,
		Src:             source,
		Dst:             tmp.Name(),
		Pwd:             pwd,
	}

	if err := client.Get(); err != nil {
		return err
	}

	success.Println("  OK")
	success.Println()

	// Load the soon-to-be-merged file
	dark.Println("Loading and parsing source")

	sourceDoc, err := pkg.Load(tmp.Name())
	if err != nil {
		return err
	}

	success.Println("  OK")
	success.Println()

	// Take current assignments and set them in the new doc
	dark.Println("Updating upstream with key/value pairs from", primary.Sprint(filename))
	dark.Println()

	sawError := false
	lastWasError := false
	counter := 0

	for _, stmt := range env.AllAssignments() {
		if !stmt.Active {
			continue
		}

		options := ast.UpsertOptions{
			SkipIfSame:     true,
			ErrorIfMissing: shared.BoolWithInverseValue(cmd.Flags(), "error-on-missing-key"),
		}

		// If the KEY does not exists in the SOURCE doc
		if sourceDoc.Get(stmt.Name) == nil {
			// Copy comments if the KEY doesn't exist in the SOURCE document
			options.Comments = stmt.CommentsSlice()

			// Try to find positioning in the statement list for the new KEY pair
			var parent ast.StatementCollection = env
			if stmt.Group != nil {
				parent = stmt.Group
			}

			idx, _ := parent.GetAssignmentIndex(stmt.Name)

			// Try to keep the position of the KEY around where it was before
			switch {
			// If we can't find any placement, put us last in the list
			case idx == -1:
				options.UpsertPlacementType = ast.UpsertLast

				// Retain the group name if its still present in the SOURCE doc
				if stmt.Group != nil && sourceDoc.HasGroup(stmt.Group.String()) {
					options.Group = stmt.Group.String()
				}

			// If we were first in the FILE doc, make sure we're first again
			case idx == 0:
				options.UpsertPlacementType = ast.UpsertFirst

				// Retain the group name if its still present in the SOURCE doc
				if stmt.Group != nil && sourceDoc.HasGroup(stmt.Group.String()) {
					options.Group = stmt.Group.String()
				}

			// If we were not first, then put us behind the key that was
			// just before us in the FILE doc
			case idx > 0:
				before := parent.Assignments()[idx-1]

				options.UpsertPlacementType = ast.UpsertAfter
				options.UpsertPlacementValue = before.Name

				if before.Group != nil && sourceDoc.HasGroup(before.Group.String()) {
					options.Group = before.Group.String()
				}
			}
		}

		changed, err := sourceDoc.Upsert(stmt, options)
		if err != nil {
			sawError = true
			lastWasError = true

			if counter > 0 {
				dark.Println()
			}

			dark.Print("  ")
			dangerEmphasis.Print(stmt.Name)
			dark.Print(" could not be set to ")
			primary.Print(stmt.Literal)
			dark.Println(" due to error:")

			danger.Println(" ", strings.Repeat(" ", len(stmt.Name)), err.Error())

			counter++

			continue
		}

		if errors := validation.ValidateSingleAssignment(env, stmt.Name, nil, []string{"file", "dir"}); len(errors) > 0 {
			sawError = true
			lastWasError = true

			if counter > 0 {
				dark.Println()
			}

			dark.Print("  ")
			dangerEmphasis.Print(stmt.Name)
			dark.Print(" could not be set to ")
			primary.Print(stmt.Literal)
			dark.Println(" due to validation error:")

			for _, errIsh := range errors {
				danger.Println(" ", strings.Repeat(" ", len(stmt.Name)), strings.TrimSpace(validation.Explain(env, errIsh, false, false)))
			}

			counter++

			continue
		}

		if changed != nil {
			counter++

			if lastWasError {
				danger.Println()
			}

			lastWasError = false

			success.Print("  ", stmt.Name)
			dark.Print(" was successfully set to ")
			primary.Println(stmt.Literal)
		}
	}

	dark.Println()

	if sawError && shared.BoolWithInverseValue(cmd.Flags(), "validate") {
		return errors.New("some fields failed validation, aborting ...")
	}

	dark.Print("Backing up ")
	primary.Print(filename)
	dark.Print(" to ")
	primary.Print(filename, ".dottie-backup")
	primary.Println()

	if err := Copy(filename, filename+".dottie-backup"); err != nil {
		danger.Println("  ERROR", err.Error())

		return err
	}

	success.Println("  OK")
	success.Println()

	dark.Println("Saving the new", primary.Sprint(filename))

	if err := pkg.Save(filename, sourceDoc); err != nil {
		danger.Println("  ERROR", err.Error())

		return err
	}

	success.Println("  OK")
	success.Println()

	success.Box("Update successfully completed")

	return nil
}

func Copy(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	info, err := srcF.Stat()
	if err != nil {
		return err
	}

	dstF, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}
