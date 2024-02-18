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
	"errors"
	"fmt"
	"slices"
	"strings"

	"go.uber.org/multierr"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

// Mapping is a user-supplied function which maps from variable names to values.
// Returns the value as a string and a bool indicating whether
// the value is present, to distinguish between an empty string
// and the absence of a value.
type Mapping func(string) (string, bool)

type Lookupper struct {
	resolver Mapping
	missing  func(string)
}

func (l Lookupper) Get(name string) expand.Variable {
	val, ok := l.resolver(name)
	if !ok {
		if name != "IFS" {
			l.missing(name)
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

func (l Lookupper) Each(cb func(name string, vr expand.Variable) bool) {
	panic("Lookupper each")
}

// SubstituteWithOptions substitute variables in the string with their values.
// It accepts additional options such as a custom function or pattern.
func Substitute(template string, mapping Mapping) (string, error, error) {
	var (
		combinedWarnings, combinedErrors error
		missing                          []string
		variables                        = ExtractVariables(template)
	)

	looker := Lookupper{
		resolver: mapping,
		missing: func(key string) {
			variable, ok := variables[key]

			// shouldn't be a lookup for anything that
			if !ok {
				panic(fmt.Errorf("unexpected missing() call during template.Substitute() for KEY [%s] - it's not in variable list?!", key))
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
		Env:     looker,
		NoUnset: false,
	}

	words, err := syntax.NewParser(syntax.Variant(syntax.LangBash)).Document(strings.NewReader(template))
	if err != nil {
		return "", nil, InvalidTemplateError{Template: template}
	}

	// Expand variables
	result, err := expand.Literal(config, words)

	target := &expand.UnsetParameterError{}
	if errors.As(err, target) {
		combinedErrors = multierr.Append(combinedErrors, &MissingRequiredError{
			Variable: target.Node.Param.Value,
			Reason:   target.Message,
		})
	}

	for _, missingKey := range missing {
		combinedWarnings = multierr.Append(combinedWarnings, fmt.Errorf("The %q variable is not set. Defaulting to a blank string.", missingKey))
	}

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
		for _, p := range w.Parts {
			switch part := p.(type) {
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
			}
		}

		return true
	})

	return variables, len(variables) > 0
}

// Split the string at the first occurrence of sep, and return the part before the separator,
// and the part after the separator.
//
// If the separator is not found, return the string itself, followed by an empty string.
func partition(str, sep string) (string, string) {
	if strings.Contains(str, sep) {
		parts := strings.SplitN(str, sep, 2)

		return parts[0], parts[1]
	}

	return str, ""
}
