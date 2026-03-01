//   Copyright 2020 The Compose Specification Authors.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package template

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"

	"github.com/jippi/dottie/pkg/tui"
	slogctx "github.com/veqryn/slog-context"
	"go.uber.org/multierr"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

// SubstituteWithOptions substitute variables in the string with their values.
// It accepts additional options such as a custom function or pattern.
func Substitute(ctx context.Context, input string, resolver Resolver, accessibleVariables AccessibleVariables) (string, error) {
	ctx = slogctx.With(ctx, slog.String("source", "template.Substitute"))

	slogctx.Debug(ctx, "template.Substitute.input", tui.StringDump("input", input))

	var combinedErrors error

	environment := EnvironmentHelper{
		Resolver:            resolver,
		MissingKeyCallback:  DefaultMissingKeyCallback(ctx, input),
		AccessibleVariables: accessibleVariables,
	}

	config := &expand.Config{
		Env: environment,
		// Any commands being tried to run will simply be treated as literals
		//
		// NOTE: the printer _will_ format the code, so that might cause some unwanted side-effects,
		//       please see https://github.com/mvdan/sh for any issues
		//
		// Example:
		//
		//  - input : $(echo hello | tee > something)
		//    output: $(echo hello | tee >something)
		//
		//  - input : ``$
		//    output: $()$
		CmdSubst: func(writer io.Writer, i *syntax.CmdSubst) error {
			start := i.Left.Offset()
			end := i.End().Offset() - 1

			writer.Write([]byte(input[start:end]))

			return nil
		},
	}

	// Parse template into Shell words
	words, err := syntax.NewParser(syntax.Variant(syntax.LangBash)).Document(strings.NewReader(input))
	if err != nil {
		return "", InvalidTemplateError{Template: input}
	}

	// Expand variables
	var result string

	err = func() (innerErr error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				// Fuzzing uncovered panics inside the third-party shell expansion library
				// for malformed arithmetic expansions. Convert that panic to a regular
				// template error so bad input fails safely without crashing the parser.
				innerErr = fmt.Errorf("template expansion panic: %v", recovered)
			}
		}()

		result, innerErr = expand.Literal(config, words)

		return innerErr
	}()
	if err != nil {
		// Inspect error and enrich it
		target := &expand.UnsetParameterError{}

		switch {
		case errors.As(err, target):
			combinedErrors = multierr.Append(combinedErrors, &MissingRequiredError{
				Variable: target.Node.Param.Value,
				Reason:   target.Message,
			})

		default:
			combinedErrors = multierr.Append(combinedErrors, InvalidTemplateError{Template: input, Wrapped: err})
		}
	}

	slogctx.Debug(ctx, "template.Substitute output", tui.StringDump("output", result))

	return result, combinedErrors
}

// ExtractVariables returns a map of all the variables defined in the specified
// composefile (dict representation) and their default value if any.
func ExtractVariables(ctx context.Context, configDict any) map[string]Variable {
	return recurseExtract(ctx, configDict)
}

func recurseExtract(ctx context.Context, value interface{}) map[string]Variable {
	results := map[string]Variable{}

	switch value := value.(type) {
	case string:
		if values, is := extractVariable(ctx, value); is {
			for _, v := range values {
				results[v.Name] = v
			}
		}

	case map[string]interface{}:
		for _, elem := range value {
			submap := recurseExtract(ctx, elem)
			for key, value := range submap {
				results[key] = value
			}
		}

	case []interface{}:
		for _, elem := range value {
			if values, is := extractVariable(ctx, elem); is {
				for _, v := range values {
					results[v.Name] = v
				}
			}
		}
	}

	return results
}

type Variable struct {
	Name          string
	DefaultValue  string
	PresenceValue string
	Required      bool
}

func extractVariable(ctx context.Context, value interface{}) ([]Variable, bool) {
	sValue, ok := value.(string)
	if !ok {
		return []Variable{}, false
	}

	var variables []Variable

	grab := func(p syntax.WordPart) string {
		switch val := p.(type) {
		case *syntax.Lit:
			return val.Value

		case *syntax.ParamExp:
			return val.Param.Value

		default:
			// Fuzzing has shown that shell parser nodes can include additional
			// word part variants; treat unknown parts as non-variable content.
			return ""
		}
	}

	slogctx.Debug(ctx, "template.extractVariable()", slog.String("sValue", sValue))

	syntax.NewParser(syntax.Variant(syntax.LangBash)).Words(strings.NewReader(sValue), func(w *syntax.Word) bool {
		for _, partInterface := range w.Parts {
			switch part := partInterface.(type) {
			case *syntax.ParamExp:
				if part.Param == nil {
					// Some malformed parameter expressions can produce ParamExp nodes
					// without a parameter name; skip those safely.
					continue
				}

				variable := Variable{
					Name: part.Param.Value,
				}

				if part.Exp != nil {
					if slices.Contains([]syntax.ParExpOperator{syntax.ErrorUnset, syntax.ErrorUnsetOrNull}, part.Exp.Op) {
						variable.Required = true
					}

					if slices.Contains([]syntax.ParExpOperator{syntax.DefaultUnsetOrNull, syntax.DefaultUnset}, part.Exp.Op) {
						if part.Exp.Word != nil && len(part.Exp.Word.Parts) > 0 {
							variable.DefaultValue = grab(part.Exp.Word.Parts[0])
						}
					}

					if slices.Contains([]syntax.ParExpOperator{syntax.AlternateUnset, syntax.AlternateUnsetOrNull}, part.Exp.Op) {
						if part.Exp.Word != nil && len(part.Exp.Word.Parts) > 0 {
							variable.PresenceValue = grab(part.Exp.Word.Parts[0])
						}
					}
				}

				variables = append(variables, variable)

			case *syntax.CmdSubst, *syntax.ProcSubst, *syntax.SglQuoted, *syntax.DblQuoted, *syntax.Lit, *syntax.ExtGlob, *syntax.ArithmExp:
				// Ignore known good-to-ignore-keywords

			default:
				// Keep extraction resilient for new/rare shell AST nodes discovered
				// by fuzzing, rather than crashing parser initialization.
				slogctx.Debug(ctx, "template.extractVariable() ignoring unsupported part", slog.String("part_type", fmt.Sprintf("%T", partInterface)))
			}
		}

		return true
	})

	return variables, len(variables) > 0
}
