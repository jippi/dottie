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
	"slices"
	"strings"

	"go.uber.org/multierr"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

// Resolver is a user-supplied function which maps from variable names to values.
// Returns the value as a string and a bool indicating whether
// the value is present, to distinguish between an empty string
// and the absence of a value.
type Resolver func(string) (string, bool)

type EnvironmentHelper struct {
	resolver           Resolver
	missingKeyCallback func(string)
}

func (helper EnvironmentHelper) Get(name string) expand.Variable {
	val, ok := helper.resolver(name)
	if !ok {
		if name != "IFS" {
			helper.missingKeyCallback(name)
		}

		return expand.Variable{
			Kind: expand.Unset,
		}
	}

	return expand.Variable{
		Str:      val,
		Exported: true,
		ReadOnly: true,
		Kind:     expand.String,
	}
}

func (l EnvironmentHelper) Each(cb func(name string, vr expand.Variable) bool) {
	panic("EnvironmentHelper.Each() should never be called")
}

// SubstituteWithOptions substitute variables in the string with their values.
// It accepts additional options such as a custom function or pattern.
func Substitute(_ context.Context, template string, resolver Resolver) (string, error, error) {
	fmt.Println("template.Substitute input:", fmt.Sprintf(">%q<", template))

	var (
		combinedWarnings, combinedErrors error
		missing                          []string
		variables                        = ExtractVariables(template)
	)

	environment := EnvironmentHelper{
		resolver: resolver,
		missingKeyCallback: func(key string) {
			variable, ok := variables[key]

			// shouldn't be a lookup for anything that
			if !ok {
				missing = append(missing, key)

				return
			}

			// Required variables are errors, so we ignore them as warnings
			if variable.Required {
				return
			}

			// If the variable has a default value, then it's not missing
			if len(variable.DefaultValue) > 0 {
				return
			}

			// If the variable has a alternate/presence value, then it's not missing
			if len(variable.PresenceValue) > 0 {
				return
			}

			missing = append(missing, key)
		},
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
			start := i.Left.Offset() - 1
			end := i.End().Offset() - 1

			writer.Write([]byte(template[start:end]))

			return nil
		},
	}

	// Parse template into Shell words
	//
	// Single quote the input to avoid shell expansions such as "~" => $HOME => env lookup
	words, err := syntax.NewParser(syntax.Variant(syntax.LangBash)).Document(strings.NewReader(template))
	if err != nil {
		return "", nil, InvalidTemplateError{Template: template}
	}

	// Expand variables
	result, err := expand.Literal(config, words)
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
			combinedErrors = multierr.Append(combinedErrors, InvalidTemplateError{Template: template, Wrapped: err})
		}
	}

	// Emit missing key warnings
	for _, missingKey := range missing {
		combinedWarnings = multierr.Append(combinedWarnings, fmt.Errorf("The [ $%s ] key is not set. Defaulting to a blank string.", missingKey))
	}

	fmt.Println("template.Substitute output:", fmt.Sprintf(">%q<", result))

	return result, combinedWarnings, combinedErrors
}

// ExtractVariables returns a map of all the variables defined in the specified
// composefile (dict representation) and their default value if any.
func ExtractVariables(configDict any) map[string]Variable {
	return recurseExtract(configDict)
}

func recurseExtract(value interface{}) map[string]Variable {
	results := map[string]Variable{}

	switch value := value.(type) {
	case string:
		if values, is := extractVariable(value); is {
			for _, v := range values {
				results[v.Name] = v
			}
		}

	case map[string]interface{}:
		for _, elem := range value {
			submap := recurseExtract(elem)
			for key, value := range submap {
				results[key] = value
			}
		}

	case []interface{}:
		for _, elem := range value {
			if values, is := extractVariable(elem); is {
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

func extractVariable(value interface{}) ([]Variable, bool) {
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
			panic(val)
		}
	}

	syntax.NewParser(syntax.Variant(syntax.LangBash)).Words(strings.NewReader(sValue), func(w *syntax.Word) bool {
		for _, partInterface := range w.Parts {
			switch part := partInterface.(type) {
			case *syntax.ParamExp:
				variable := Variable{
					Name: part.Param.Value,
				}

				if part.Exp != nil {
					if slices.Contains([]syntax.ParExpOperator{syntax.ErrorUnset, syntax.ErrorUnsetOrNull}, part.Exp.Op) {
						variable.Required = true
					}

					if slices.Contains([]syntax.ParExpOperator{syntax.DefaultUnsetOrNull, syntax.DefaultUnset}, part.Exp.Op) {
						variable.DefaultValue = grab(part.Exp.Word.Parts[0])
					}

					if slices.Contains([]syntax.ParExpOperator{syntax.AlternateUnset, syntax.AlternateUnsetOrNull}, part.Exp.Op) {
						variable.PresenceValue = grab(part.Exp.Word.Parts[0])
					}
				}

				variables = append(variables, variable)

			case *syntax.CmdSubst, *syntax.SglQuoted, *syntax.DblQuoted, *syntax.Lit, *syntax.ExtGlob:
				// Ignore known good-to-ignore-keywords

			default:
				panic(fmt.Errorf("unexpected type: %T", partInterface))
			}
		}

		return true
	})

	return variables, len(variables) > 0
}
