package main

import (
	"os"
	"strings"

	goversion "github.com/caarlos0/go-version"
	"github.com/davecgh/go-spew/spew"
	"github.com/jippi/dottie/cmd/disable"
	"github.com/jippi/dottie/cmd/enable"
	"github.com/jippi/dottie/cmd/fmt"
	"github.com/jippi/dottie/cmd/groups"
	"github.com/jippi/dottie/cmd/json"
	print_cmd "github.com/jippi/dottie/cmd/print"
	"github.com/jippi/dottie/cmd/set"
	"github.com/jippi/dottie/cmd/update"
	"github.com/jippi/dottie/cmd/validate"
	"github.com/jippi/dottie/cmd/value"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals
var (
	commit    = ""
	date      = ""
	treeState = ""
	version   = ""
)

const globalOptionsTemplate = `{{if .VisibleFlags}}
GLOBAL OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}{{if .Copyright}}
{{end}}
`

var (
	env      *ast.Document
	settings *render.Settings
)

func main() {
	__configureSpew()

	cobra.EnableCommandSorting = false

	root := &cobra.Command{
		Use:           "dottie",
		Short:         "Simplify working with .env files",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       buildVersion().String(),
	}

	root.AddGroup(&cobra.Group{ID: "manipulate", Title: "Manipulation Commands"})
	root.AddGroup(&cobra.Group{ID: "output", Title: "Output Commands"})

	root.AddCommand(set.Command())
	root.AddCommand(fmt.Command)
	root.AddCommand(enable.Command)
	root.AddCommand(disable.Command)
	root.AddCommand(update.Command())

	root.AddCommand(print_cmd.Command())
	root.AddCommand(value.Command)
	root.AddCommand(validate.Command())
	root.AddCommand(groups.Command)
	root.AddCommand(json.Command)

	root.PersistentFlags().StringP("file", "f", ".env", "Load this file")

	if c, err := root.ExecuteC(); err != nil {
		tui.Theme.Danger.StderrPrinter(tui.WithEmphasis(true)).Println(c.ErrPrefix(), err.Error())
		tui.Theme.Info.StderrPrinter().Printfln("Run '%v --help' for usage.\n", c.CommandPath())

		os.Exit(1)
	}
}

func __configureSpew() {
	spew.Config.DisablePointerMethods = true
	spew.Config.DisableMethods = true
}

func indent(in string) string {
	return strings.TrimSpace(strings.Join(strings.Split(in, "\n"), "\n   "))
}

func buildVersion() goversion.Info {
	return goversion.GetVersionInfo(
		// goversion.WithAppDetails("dottie", "Making .env file management easy", "https://github.com/jippi/dottie"),
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
