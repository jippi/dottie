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
	__load()

	root := &cobra.Command{
		Use:           "dottie",
		Short:         "Simplify working with .env files",
		SilenceErrors: true,
		Version:       buildVersion().String(),
	}

	root.AddCommand(disable.Command)
	root.AddCommand(enable.Command)
	root.AddCommand(fmt.Command)
	root.AddCommand(groups.Command)
	root.AddCommand(json.Command)
	root.AddCommand(print_cmd.Command())
	root.AddCommand(set.Command())
	root.AddCommand(update.Command)
	root.AddCommand(validate.Command)
	root.AddCommand(value.Command)

	root.PersistentFlags().String("file", ".env", "Load this file")

	if c, err := root.ExecuteC(); err != nil {
		tui.Theme.Danger.StderrPrinter().Println(c.ErrPrefix(), err.Error())
		tui.Theme.Info.StderrPrinter().Printfln("Run '%v --help' for usage.\n", c.CommandPath())

		os.Exit(1)
	}
}

func __load() {
	spew.Config.DisablePointerMethods = true
	spew.Config.DisableMethods = true

	spew.Dump()
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
