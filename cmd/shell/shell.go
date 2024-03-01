package shell

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jippi/dottie/pkg/console"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "shell",
		Short:   "Dottie shell",
		GroupID: "manipulate",
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(console.NewModel(cmd))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
