package main

import (
	"context"
	"io"
	"log"
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
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/urfave/cli/v3"
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

	origHelpPrinterCustom := cli.HelpPrinterCustom
	defer func() {
		cli.HelpPrinterCustom = origHelpPrinterCustom
	}()

	app := &cli.Command{
		Name:                       "dottie",
		Version:                    indent(buildVersion().String()),
		Suggest:                    true,
		EnableShellCompletion:      true,
		ShellCompletionCommandName: "completions",
		Flags:                      shared.GlobalFlags,
		Commands: []*cli.Command{
			disable.Command,
			enable.Command,
			fmt.Command,
			groups.Command,
			json.Command,
			print_cmd.Command,
			set.Command,
			update.Command,
			validate.Command,
			value.Command,
		},
	}

	cli.HelpPrinterCustom = func(out io.Writer, templ string, data interface{}, customFuncs map[string]interface{}) {
		origHelpPrinterCustom(out, templ, data, customFuncs)

		if data != app {
			origHelpPrinterCustom(app.Writer, globalOptionsTemplate, app, nil)
		}
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
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
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}

			if treeState != "" {
				i.GitTreeState = treeState
			}

			if date != "" {
				i.BuildDate = date
			}

			if version != "" {
				i.GitVersion = version
			}
		},
	)
}
