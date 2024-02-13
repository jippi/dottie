package shared

import (
	"strings"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

type Completer struct {
	options         []render.SettingsOption
	handlers        []render.Handler
	suffix          string
	suffixIsLiteral bool
}

type CobraCompleter func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func NewCompleter() *Completer {
	return (&Completer{}).
		WithSettings(render.WithOutputType(render.CompletionKeyOnly))
}

func (c *Completer) WithKeySuffix(suffix string) *Completer {
	c.suffix = suffix

	return c
}

func (c *Completer) WithSuffixIsLiteral(b bool) *Completer {
	c.suffixIsLiteral = b

	return c
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
		filename := cmd.Flag("file").Value.String()

		doc, warn, err := pkg.Load(filename)
		if warn != nil {
			tui.Theme.Warning.StderrPrinter().Println(warn)
		}
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		c.handlers = append(
			c.handlers,
			render.ExcludeComments,
			render.ExcludeHiddenViaAnnotation,
			render.RetainKeyPrefix(toComplete),
		)

		lines := render.
			NewUnfilteredRenderer(render.NewSettings(c.options...), c.handlers...).
			Statement(doc).
			Lines()

		if c.suffixIsLiteral && strings.HasSuffix(toComplete, "=") {
			key := strings.TrimSuffix(toComplete, "=")

			if assignment := doc.Get(key); assignment != nil {
				return []string{assignment.Name + "=" + assignment.Literal}, cobra.ShellCompDirectiveDefault
			}
		}

		switch len(lines) {
		case 0:
			return lines, cobra.ShellCompDirectiveNoSpace

		case 1:
			if c.suffixIsLiteral {
				// The key is the first part of a line when split by "\t".
				//
				// The "\t" is separator between the value to complete, and its documentation
				key := strings.Split(lines[0], "\t")[0]

				assignment := doc.Get(key)

				if assignment != nil {
					return []string{assignment.Name + "=" + assignment.Literal}, cobra.ShellCompDirectiveDefault
				}
			}

			return []string{lines[0] + c.suffix}, cobra.ShellCompDirectiveNoSpace

		default:
			return lines, cobra.ShellCompDirectiveNoSpace
		}
	}
}
