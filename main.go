package main

import (
	"context"
	"os"

	"github.com/jippi/dottie/cmd"
)

func main() {
	_, err := cmd.RunCommand(context.Background(), os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(1)
	}
}
