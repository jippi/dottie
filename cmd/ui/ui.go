package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jippi/dottie/pkg"
	pkgui "github.com/jippi/dottie/pkg/ui"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ui",
		Short:   "Dottie ui",
		GroupID: "manipulate",
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	document, err := pkg.Load(cmd.Context(), filename)
	if err != nil {
		return err
	}

	p := tea.NewProgram(pkgui.NewModel(cmd.Context(), document), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err = p.Run()

	return err
}
