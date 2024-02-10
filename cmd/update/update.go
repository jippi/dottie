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
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the .env file from a source",
		RunE:  runE,
	}

	cmd.Flags().String("source", "", "URL or local file path to the upstream source file. This will take precedence over any [@dottie/source] annotation in the file")

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

	mergedEnv, err := pkg.Load(tmp.Name())
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

	for _, stmt := range env.Assignments() {
		if !stmt.Active {
			continue
		}

		changed, err := mergedEnv.Upsert(stmt, ast.UpsertOptions{SkipIfSame: true, ErrorIfMissing: true})
		if err != nil {
			danger.Println("  ERROR", err.Error())

			continue
		}

		if errors := validation.ValidateSingleAssignment(env, stmt.Name, nil, []string{"file", "dir"}); len(errors) > 0 {
			sawError = true
			lastWasError = true

			dark.Println()
			dark.Print("  ")
			dangerEmphasis.Print(stmt.Name)
			dark.Print(" could not be set to ")
			primary.Print(stmt.Literal)
			dark.Println(" due to validation error:")

			for _, errIsh := range errors {
				danger.Println(" ", strings.Repeat(" ", len(stmt.Name)), strings.TrimSpace(validation.Explain(env, errIsh, false, false)))
			}

			continue
		}

		if changed != nil {
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

	if sawError {
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

	if err := pkg.Save(filename, mergedEnv); err != nil {
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
