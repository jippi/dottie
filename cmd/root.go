package cmd

import (
	"context"
	"io"
	"strings"

	goversion "github.com/caarlos0/go-version"
	disable_cmd "github.com/jippi/dottie/cmd/disable"
	enable_cmd "github.com/jippi/dottie/cmd/enable"
	exec_cmd "github.com/jippi/dottie/cmd/exec"
	fmt_cmd "github.com/jippi/dottie/cmd/fmt"
	groups_cmd "github.com/jippi/dottie/cmd/groups"
	json_cmd "github.com/jippi/dottie/cmd/json"
	print_cmd "github.com/jippi/dottie/cmd/print"
	set_cmd "github.com/jippi/dottie/cmd/set"
	shell_cmd "github.com/jippi/dottie/cmd/shell"
	template_cmd "github.com/jippi/dottie/cmd/template"
	update_cmd "github.com/jippi/dottie/cmd/update"
	validate_cmd "github.com/jippi/dottie/cmd/validate"
	value_cmd "github.com/jippi/dottie/cmd/value"
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

	root.AddCommand(set_cmd.New())
	root.AddCommand(update_cmd.New())
	root.AddCommand(fmt_cmd.New())
	root.AddCommand(disable_cmd.New())
	root.AddCommand(enable_cmd.New())
	root.AddCommand(exec_cmd.New())
	root.AddCommand(shell_cmd.New())

	root.AddCommand(print_cmd.New())
	root.AddCommand(validate_cmd.New())
	root.AddCommand(value_cmd.New())
	root.AddCommand(groups_cmd.New())
	root.AddCommand(json_cmd.New())
	root.AddCommand(template_cmd.New())

	return root
}

func RunCommand(ctx context.Context, args []string, stdout, stderr io.Writer) (*cobra.Command, error) {
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
