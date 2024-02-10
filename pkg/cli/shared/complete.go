package shared

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jippi/dottie/pkg/render"
	"github.com/spf13/cobra"
)

type Completer struct {
	options  []render.SettingsOption
	handlers []render.Handler
}

type CobraCompleter func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func NewCompleter() *Completer {
	return (&Completer{}).
		WithSettings(render.WithOutputType(render.CompletionKeyOnly))
}

func (c *Completer) WithHandlers(handlers ...render.Handler) *Completer {
	c.handlers = append(c.handlers, handlers...)

	return c
}

func (c *Completer) WithSettings(options ...render.SettingsOption) *Completer {
	c.options = append(c.options, options...)

	return c
}

func (c *Completer) Get() CobraCompleter {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		env, settings, err := Setup(cmd.Flags())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		spew.Dump(toComplete)

		settings.Apply(c.options...)
		settings.Apply(render.WithFilterKeyPrefix(toComplete))

		lines := render.
			NewUnfilteredRenderer(*settings, c.handlers...).
			Statement(env).
			Lines()

		switch len(lines) {
		case 0:
			return lines, cobra.ShellCompDirectiveNoSpace

		case 1:
			return []string{lines[0] + "="}, cobra.ShellCompDirectiveNoSpace

		default:
			return lines, cobra.ShellCompDirectiveNoSpace
		}
	}
}
