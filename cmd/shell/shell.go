package shell

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jippi/dottie/pkg/tui"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func New() *cobra.Command {
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
	app := console.New("dottie")
	app.Shell().AcceptMultiline = func(line []rune) (accept bool) { return true }

	menu := app.ActiveMenu()
	menu.SetErr(cmd.ErrOrStderr())
	menu.SetOut(cmd.OutOrStdout())
	menu.AddHistorySourceFile("default", ".dottie.hist")
	menu.AddInterrupt(io.EOF, exitCtrlD)
	menu.SetCommands(mainMenuCommands(cmd.Root(), app))
	menu.ErrorHandler = func(err error) error {
		tui.StderrFromContext(cmd.Context()).Danger().Println(err)

		return nil
	}

	setupPrompt(menu)

	return app.Start(cmd.Context())
}

func mainMenuCommands(rootCmd *cobra.Command, _ *console.Console) console.Commands {
	// Commander
	processCommand := func(cmd *cobra.Command) {
		carapaceCompleter := carapace.Gen(cmd)
		flagMap := make(carapace.ActionMap)

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Name == "file" || strings.Contains(f.Usage, "file") {
				flagMap[f.Name] = carapace.ActionFiles()
			}
		})

		carapaceCompleter.FlagCompletion(flagMap)
	}

	return func() *cobra.Command {
		processCommand(rootCmd)

		for _, cmd := range rootCmd.Commands() {
			processCommand(cmd)
		}

		rootCmd.InitDefaultHelpCmd()
		rootCmd.CompletionOptions.DisableDefaultCmd = true
		// rootCmd.DisableFlagsInUseLine = true

		return rootCmd
	}
}

// exitCtrlD is a custom interrupt handler to use when the shell
// readline receives an io.EOF error, which is returned with CtrlD.
func exitCtrlD(c *console.Console) {
	os.Exit(0)
}

// setupPrompt is a function which sets up the prompts for the main menu.
func setupPrompt(m *console.Menu) {
	prompt := m.Prompt()

	prompt.Primary = func() string {
		prompt := "\x1b[33mdottie\x1b[0m in \x1b[34m%s\x1b[0m\n> "
		wd, _ := os.Getwd()

		dir, err := filepath.Rel(os.Getenv("HOME"), wd)
		if err != nil {
			dir = filepath.Base(wd)
		}

		return fmt.Sprintf(prompt, dir)
	}

	prompt.Secondary = func() string { return ">" }
	prompt.Right = func() string {
		return "\x1b[1;30m" + time.Now().Format("03:04:05.000") + "\x1b[0m"
	}

	prompt.Transient = func() string { return "\x1b[1;30m" + ">> " + "\x1b[0m" }
}
