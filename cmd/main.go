package main

import (
	"context"
	"io"
	"log"
	"os"

	"dotfedi/pkg/ast"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli/v3"
)

const globalOptionsTemplate = `{{if .VisibleFlags}}
GLOBAL OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}{{if .Copyright}}
{{end}}
`

var (
	env      *ast.File
	settings *ast.RenderSettings
)

func main() {
	__load()

	// spew.Config.DisableMethods = true
	// spew.Config.DisablePointerMethods = true

	origHelpPrinterCustom := cli.HelpPrinterCustom
	defer func() {
		cli.HelpPrinterCustom = origHelpPrinterCustom
	}()

	app := &cli.Command{
		Flags: globalFlags,
		Commands: []*cli.Command{
			disableCommand,
			enableCommand,
			groupsCommand,
			printCommand,
			setCommand,
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
	spew.Dump()
}
