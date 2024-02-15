package main

import (
	"os"

	"github.com/jippi/dottie/cmd"
	"github.com/jippi/dottie/pkg/tui"
)

func main() {
	root := cmd.NewCommand()
	if c, err := root.ExecuteC(); err != nil {
		tui.Theme.Danger.BuffPrinter(root.ErrOrStderr(), tui.WithEmphasis(true)).Printfln("%s %+v", c.ErrPrefix(), err)
		tui.Theme.Info.BuffPrinter(root.ErrOrStderr()).Printfln("Run '%v --help' for usage.\n", c.CommandPath())

		os.Exit(1)
	}
}
