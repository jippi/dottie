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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update the .env file from a source",
		GroupID: "manipulate",
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}

	cmd.Flags().String("source", "", "URL or local file path to the upstream source file. This will take precedence over any [@dottie/source] annotation in the file")
	cmd.Flags().StringSlice("ignore-rule", []string{}, "Ignore this validation rule (e.g. 'dir')")
	cmd.Flags().StringSlice("exclude-key-prefix", []string{}, "Ignore these KEY prefixes")

	shared.BoolWithInverse(cmd, "error-on-missing-key", true, "Error if a KEY in FILE is missing from SOURCE", "Add KEY to FILE if missing from SOURCE")
	shared.BoolWithInverse(cmd, "validate", true, "Validation errors will abort the update", "Validation errors will be printed but will not fail the update")
	shared.BoolWithInverse(cmd, "save", true, "Save the document after processing", "Do not save the document after processing")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	stdout, stderr := tui.WritersFromContext(cmd.Context())

	dark := stdout.NoColor()
	info := stdout.Info()
	danger := stdout.Danger()
	success := stdout.Success()
	primary := stdout.Primary()

	info.Box("Starting update of " + filename + " from upstream")
	info.Println()

	dark.Println("Looking for upstream source")

	oldDocument, err := pkg.Load(cmd.Context(), filename)
	if err != nil {
		return err
	}

	source, _ := cmd.Flags().GetString("source")
	if len(source) == 0 {
		source, err = oldDocument.GetConfig("dottie/source")
		if err != nil {
			return err
		}

		success.Println("  Found source via [dottie/source] annotation in file", primary.Sprint(filename))
	} else {
		success.Println("  Found source via CLI flag")
	}

	dark.Println()

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

	newDocument, err := pkg.Load(cmd.Context(), tmp.Name())
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

	var selectors []ast.Selector

	selectors = append(selectors, ast.ExcludeDisabledAssignments)

	if slice := shared.StringSliceFlag(cmd.Flags(), "exclude-key-prefix"); len(slice) > 0 {
		for _, prefix := range slice {
			selectors = append(selectors, ast.ExcludeKeyPrefix(prefix))
		}
	}

	for _, oldStatement := range oldDocument.AllAssignments(selectors...) {
		if !oldStatement.Enabled {
			continue
		}

		upserter, err := upsert.New(
			newDocument,
			upsert.EnableSetting(upsert.SkipIfSame),
			upsert.EnableSetting(upsert.SkipIfEmpty),
			upsert.EnableSetting(upsert.SkipIfSet),
			upsert.EnableSettingIf(upsert.ErrorIfMissing, shared.BoolWithInverseValue(cmd.Flags(), "error-on-missing-key")),
			upsert.WithSkipValidationRule(shared.StringSliceFlag(cmd.Flags(), "ignore-rule")...),
		)
		if err != nil {
			return err
		}

		// If the KEY does *NOT* exists in the SOURCE doc
		if newDocument.Get(oldStatement.Name) == nil {
			// Try to find positioning in the statement list for the new KEY pair
			var parent ast.StatementCollection = oldDocument

			upserter.ApplyOptions(upsert.EnableSetting(upsert.UpdateComments))

			if oldStatement.Group != nil {
				parent = oldStatement.Group
			}

			idx, _ := parent.GetAssignmentIndex(oldStatement.Name)

			// Try to keep the position of the KEY around where it was before
			switch {
			// If we can't find any placement, put us last in the list
			case idx == -1:
				upserter.ApplyOptions(upsert.WithPlacement(upsert.AddLast))

				// Retain the group name if its still present in the SOURCE doc
				if oldStatement.Group != nil && newDocument.HasGroup(oldStatement.Group.String()) {
					upserter.ApplyOptions(upsert.WithGroup(oldStatement.Group.String()))
				}

			// If we were first in the FILE doc, make sure we're first again
			case idx == 0:
				upserter.ApplyOptions(upsert.WithPlacement(upsert.AddFirst))

				// Retain the group name if its still present in the SOURCE doc
				if oldStatement.Group != nil && newDocument.HasGroup(oldStatement.Group.String()) {
					upserter.ApplyOptions(upsert.WithGroup(oldStatement.Group.String()))
				}

			// If we were not first, then put us behind the key that was
			// just before us in the FILE doc
			case idx > 0:
				// Search previous keys to find a possible group to add the new KEY to
				for ; idx > 0; idx-- {
					before := parent.Assignments()[idx-1]

					if err := upserter.ApplyOptions(upsert.WithPlacementRelativeToKey(upsert.AddAfterKey, before.Name)); err != nil {
						return err
					}

					if before.Group != nil && newDocument.HasGroup(before.Group.String()) {
						upserter.ApplyOptions(upsert.WithGroup(before.Group.String()))

						break
					}

					if x := newDocument.Get(before.Name); x != nil && x.Group != nil {
						upserter.ApplyOptions(upsert.WithGroup(x.Group.String()))

						break
					}
				}
			}
		}

		var skippedStatementWarning upsert.SkippedStatementError

		changed, err := upserter.Upsert(cmd.Context(), oldStatement)

		switch {
		case errors.As(err, &skippedStatementWarning):
			stderr.Warning().Print("WARNING: Key [", oldStatement.Name, "] was skipped: ")
			stderr.Warning().Println(skippedStatementWarning.Reason)

		case err != nil:
			sawError = true
			lastWasError = true

			if counter > 0 {
				dark.Println()
			}

			danger.Print(validation.Explain(cmd.Context(), newDocument, err, changed, false, true))

			counter++

			continue
		}

		if changed != nil {
			counter++

			if lastWasError {
				danger.Println()
			}

			lastWasError = false

			success.Print("  ", oldStatement.Name)
			dark.Print(" was successfully set to ")
			primary.Println(oldStatement.Literal)
		}
	}

	stdout.NoColor().Println()

	if sawError && shared.BoolWithInverseValue(cmd.Flags(), "validate") {
		stdout.NoColor().Println()

		return errors.New("some fields failed validation, aborting ...")
	}

	if !shared.BoolWithInverseValue(cmd.Flags(), "save") {
		stdout.Warning().Println("[--no-save] was provided, not saving file")

		return nil
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

	if err := pkg.Save(cmd.Context(), filename, newDocument); err != nil {
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

func indent(in string, width int) string {
	return strings.Repeat(" ", width) + strings.TrimSpace(strings.Join(strings.Split(in, "\n"), "\n"+strings.Repeat(" ", width)))
}
