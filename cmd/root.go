package cmd

import (
	"context"
	"io"
	"strings"

	goversion "github.com/caarlos0/go-version"
	"github.com/jippi/dottie/cmd/disable"
	"github.com/jippi/dottie/cmd/enable"
	"github.com/jippi/dottie/cmd/exec"
	"github.com/jippi/dottie/cmd/fmt"
	"github.com/jippi/dottie/cmd/groups"
	"github.com/jippi/dottie/cmd/json"
	print_cmd "github.com/jippi/dottie/cmd/print"
	"github.com/jippi/dottie/cmd/set"
	"github.com/jippi/dottie/cmd/shell"
	"github.com/jippi/dottie/cmd/update"
	"github.com/jippi/dottie/cmd/validate"
	"github.com/jippi/dottie/cmd/value"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals
var (
	commit    = "UNSET"
	date      = "UNSET"
	treeState = "UNSET"
	version   = "dev"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:     "dottie",
		Short:   "Simplify working with .env files",
		Version: buildVersion().String(),
	}

	root.AddGroup(&cobra.Group{ID: "manipulate", Title: "Manipulation Commands"})
	root.AddGroup(&cobra.Group{ID: "output", Title: "Output Commands"})

	root.AddCommand(set.NewCommand())
	root.AddCommand(exec.NewCommand())
	root.AddCommand(fmt.NewCommand())
	root.AddCommand(enable.NewCommand())
	root.AddCommand(disable.NewCommand())
	root.AddCommand(update.NewCommand())
	root.AddCommand(shell.NewCommand())

	root.AddCommand(print_cmd.NewCommand())
	root.AddCommand(value.NewCommand())
	root.AddCommand(validate.NewCommand())
	root.AddCommand(groups.NewCommand())
	root.AddCommand(json.NewCommand())

	return root
}

func RunCommand(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) (*cobra.Command, error) {
	root := NewRootCommand()
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetArgs(args)
	root.SetContext(ctx)
	root.SetErr(stderr)
	root.SetOut(stdout)
	root.PersistentFlags().StringP("file", "f", ".env", "Load this file")
	root.SetVersionTemplate(`{{ .Version }}`)

	command, err := root.ExecuteC()
	if err != nil {
		stderr := tui.WriterFromContext(ctx, tui.Stderr)
		stderr.Danger().Copy(tui.WithEmphasis(true)).Printfln("%s %+v", command.ErrPrefix(), err)
		stderr.Info().Printfln("Run '%v --help' for usage.", command.CommandPath())
	}

	return command, err
}

func indent(in string) string {
	return strings.TrimSpace(strings.Join(strings.Split(in, "\n"), "\n   "))
}

func buildVersion() goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails("dottie", "Making .env file management easy", "https://github.com/jippi/dottie"),
		func(versionInfo *goversion.Info) {
			if commit != "" {
				versionInfo.GitCommit = commit
			}

			if treeState != "" {
				versionInfo.GitTreeState = treeState
			}

			if date != "" {
				versionInfo.BuildDate = date
			}

			if version != "" {
				versionInfo.GitVersion = version
			}
		},
	)
}
