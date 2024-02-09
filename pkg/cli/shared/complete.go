package shared

import (
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

		settings.Apply(render.WithFilterKeyPrefix(toComplete))
		settings.Apply(c.options...)

		return render.
				NewUnfilteredRenderer(*settings, c.handlers...).
				Statement(env).
				Lines(),
			cobra.ShellCompDirectiveDefault
	}
}
