package update

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/go-getter"
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "update",
	Short: "Update the .env file from a source",
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := cmd.Flag("file").Value.String()

		env, err := pkg.Load(filename)
		if err != nil {
			return err
		}

		fmt.Print("Finding source")
		source, err := env.GetConfig("dottie/source")
		if err != nil {
			return err
		}
		fmt.Println(" ✅")

		fmt.Print("Grabbing .env.docker from [", source, "]")

		if _, err := os.Stat(".env.source"); errors.Is(err, os.ErrNotExist) {
			tmp, err := os.OpenFile(".env.source", os.O_RDWR|os.O_CREATE, 0o666)
			if err != nil {
				return err
			}

			// Get the pwd
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Error getting wd: %s", err)
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
		}
		fmt.Println(" ✅")

		// Copy source to "new"
		fmt.Print("Copying .env.source into .env.merged")
		if err := Copy(".env.source", ".env.merged"); err != nil {
			return err
		}
		fmt.Println(" ✅")

		// Load the soon-to-be-merged file
		fmt.Print("Loading and parsing .env.merged")
		mergedEnv, err := pkg.Load(".env.merged")
		if err != nil {
			return err
		}
		fmt.Println(" ✅")

		// Take current assignments and set them in the new doc
		fmt.Println("Updating .env.merged with key/value pairs from .env")
		for _, stmt := range env.Assignments() {
			if !stmt.Active {
				continue
			}

			changed, err := mergedEnv.Upsert(stmt, ast.UpsertOptions{SkipIfSame: true, ErrorIfMissing: true})
			if err != nil {
				fmt.Println("  ❌", err.Error())

				continue
			}

			if changed != nil {
				fmt.Println("  ✅", fmt.Sprintf("Key [%s] was successfully updated", stmt.Name))
			}
		}

		fmt.Println()
		fmt.Println("Saving .env.merged")

		return pkg.Save(".env.merged", mergedEnv)
	},
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