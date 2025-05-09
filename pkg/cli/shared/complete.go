package shared

import (
	"strings"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/render"
	"github.com/spf13/cobra"
)

type Completer struct {
	options         []render.SettingsOption
	selectors       []ast.Selector
	suffix          string
	suffixIsLiteral bool
}

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

func (c *Completer) WithSelectors(selectors ...ast.Selector) *Completer {
	c.selectors = append(c.selectors, selectors...)

	return c
}

func (c *Completer) WithSettings(options ...render.SettingsOption) *Completer {
	c.options = append(c.options, options...)

	return c
}

func (c *Completer) Get() cobra.CompletionFunc {
	c.selectors = append(
		c.selectors,
		ast.ExcludeComments,
		ast.ExcludeHiddenViaAnnotation,
	)

	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		filename := cmd.Flag("file").Value.String()

		doc, err := pkg.Load(cmd.Context(), filename)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		selectors := c.selectors

		if len(toComplete) > 0 {
			selectors = append(selectors, ast.RetainKeyPrefix(toComplete))
		}

		lines := render.
			NewUnfilteredRenderer(render.NewSettings(c.options...), render.NewAstSelectorHandler(selectors...)).
			Statement(cmd.Context(), doc).
			Lines()

		if c.suffixIsLiteral && strings.HasSuffix(toComplete, "=") {
			key := strings.TrimSuffix(toComplete, "=")

			if assignment := doc.Get(key); assignment != nil {
				return []string{assignment.Name + "=" + assignment.GetSafeLiteral()}, cobra.ShellCompDirectiveNoSpace
			}
		}

		switch len(lines) {
		case 0:
			return lines, cobra.ShellCompDirectiveNoFileComp

		case 1:
			if c.suffixIsLiteral {
				// The key is the first part of a line when split by "\t".
				//
				// The "\t" is separator between the value to complete, and its documentation
				key := strings.Split(lines[0], "\t")[0]

				assignment := doc.Get(key)

				if assignment != nil {
					return []string{assignment.Name + "=" + assignment.GetSafeLiteral()}, cobra.ShellCompDirectiveNoFileComp
				}
			}

			return []string{lines[0] + c.suffix}, cobra.ShellCompDirectiveNoFileComp

		default:
			return lines, cobra.ShellCompDirectiveNoFileComp
		}
	}
}
