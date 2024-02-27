package main

import (
	"context"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/jippi/dottie/cmd"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

func main() {
	spew.Config.DisablePointerMethods = false
	spew.Config.DisableMethods = false
	cobra.EnableCommandSorting = false

	ctx := tui.NewContext(context.Background(), os.Stdout, os.Stderr)

	_, err := cmd.RunCommand(ctx, os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(1)
	}
}
