package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	goversion "github.com/caarlos0/go-version"
	"github.com/davecgh/go-spew/spew"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/render"
	"github.com/urfave/cli/v3"
)

// nolint: gochecknoglobals
var (
	builtBy   = ""
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
		Flags:                      globalFlags,
		Commands: []*cli.Command{
			disableCommand,
			enableCommand,
			groupsCommand,
			formatCmd,
			jsonCommand,
			printCommand,
			setCommand,
			updateCommand,
			validateCommand,
			valueCommand,
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

			if builtBy != "" {
				i.BuiltBy = builtBy
			}
		},
	)
}
