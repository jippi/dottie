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
	"github.com/jippi/dottie/pkg/ast/upsert"
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

	originalEnv, err := pkg.Load(filename)
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
		source, err = originalEnv.GetConfig("dottie/source")
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

	for _, originalStatement := range originalEnv.AllAssignments() {
		if !originalStatement.Enabled {
			continue
		}

		upserter, err := upsert.New(
			sourceDoc,
			upsert.WithSetting(upsert.SkipIfSame),
			upsert.WithSettingIf(upsert.ErrorIfMissing, shared.BoolWithInverseValue(cmd.Flags(), "error-on-missing-key")),
		)
		if err != nil {
			return err
		}

		// If the KEY does *NOT* exists in the SOURCE doc
		if sourceDoc.Get(originalStatement.Name) == nil {
			// Try to find positioning in the statement list for the new KEY pair
			var parent ast.StatementCollection = originalEnv

			if originalStatement.Group != nil {
				parent = originalStatement.Group
			}

			idx, _ := parent.GetAssignmentIndex(originalStatement.Name)

			// Try to keep the position of the KEY around where it was before
			switch {
			// If we can't find any placement, put us last in the list
			case idx == -1:
				upserter.ApplyOptions(upsert.WithPlacement(upsert.AddLast))

				// Retain the group name if its still present in the SOURCE doc
				if originalStatement.Group != nil && sourceDoc.HasGroup(originalStatement.Group.String()) {
					upserter.ApplyOptions(upsert.WithGroup(originalStatement.Group.String()))
				}

			// If we were first in the FILE doc, make sure we're first again
			case idx == 0:
				upserter.ApplyOptions(upsert.WithPlacement(upsert.AddFirst))

				// Retain the group name if its still present in the SOURCE doc
				if originalStatement.Group != nil && sourceDoc.HasGroup(originalStatement.Group.String()) {
					upserter.ApplyOptions(upsert.WithGroup(originalStatement.Group.String()))
				}

			// If we were not first, then put us behind the key that was
			// just before us in the FILE doc
			case idx > 0:
				before := parent.Assignments()[idx-1]

				if err := upserter.ApplyOptions(upsert.WithPlacementRelativeToKey(upsert.AddAfterKey, before.Name)); err != nil {
					return err
				}

				if before.Group != nil && sourceDoc.HasGroup(before.Group.String()) {
					upserter.ApplyOptions(upsert.WithGroup(before.Group.String()))
				}
			}
		}

		changed, warn, err := upserter.Upsert(originalStatement)
		if warn != nil {
			tui.Theme.Warning.StderrPrinter().Println(warn)
		}

		if err != nil {
			sawError = true
			lastWasError = true

			if counter > 0 {
				dark.Println()
			}

			dark.Print("  ")
			dangerEmphasis.Print(originalStatement.Name)
			dark.Print(" could not be set to ")
			primary.Print(originalStatement.Literal)
			dark.Println(" due to error:")

			danger.Println(" ", strings.Repeat(" ", len(originalStatement.Name)), err.Error())

			counter++

			continue
		}

		if errors := validation.ValidateSingleAssignment(originalEnv, originalStatement, nil, []string{"file", "dir"}); len(errors) > 0 {
			sawError = true
			lastWasError = true

			if counter > 0 {
				dark.Println()
			}

			dark.Print("  ")
			dangerEmphasis.Print(originalStatement.Name)
			dark.Print(" could not be set to ")
			primary.Print(originalStatement.Literal)
			dark.Println(" due to validation error:")

			for _, errIsh := range errors {
				danger.Println(" ", strings.Repeat(" ", len(originalStatement.Name)), strings.TrimSpace(validation.Explain(originalEnv, errIsh, errIsh, false, false)))
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

			success.Print("  ", originalStatement.Name)
			dark.Print(" was successfully set to ")
			primary.Println(originalStatement.Literal)
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
