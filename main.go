package main

import (
	"context"
	"os"

	"github.com/jippi/dottie/cmd"
	"github.com/jippi/dottie/pkg/tui"
)

func main() {
	ctx := tui.NewContext(context.Background(), os.Stdout, os.Stderr)

	_, err := cmd.RunCommand(ctx, os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(1)
	}
}
