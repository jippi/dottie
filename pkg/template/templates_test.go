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

package template_test

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	templatepkg "github.com/jippi/dottie/pkg/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaults = map[string]string{
	"FOO":  "first",
	"BAR":  "",
	"JSON": `{"json":2}`,
}

var MissingRequiredError = &templatepkg.MissingRequiredError{}

func defaultMapping(name string) (string, bool) {
	val, ok := defaults[name]

	return val, ok
}

func TestEscaped(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("$${foo}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "${foo}", result)
}

func TestSubstituteNoMatch(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("foo", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "foo", result)
}

func TestUnescaped(t *testing.T) {
	t.Parallel()

	templates := []string{
		"a $ string",
		"^REGEX$",
		"$}",
		"$",
	}

	for _, expected := range templates {
		actual, warn, err := templatepkg.Substitute(expected, defaultMapping)
		assert.NoError(t, warn)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestInvalid(t *testing.T) {
	t.Parallel()

	invalidTemplates := []string{
		// "${",
		// "${}",
		"${ }",
		"${ foo}",
		"${foo }",
		"${foo!}",
	}

	for i, tt := range invalidTemplates {
		tt := tt

		t.Run(fmt.Sprintf("TestInvalid %d", i), func(t *testing.T) {
			t.Parallel()

			_, _, err := templatepkg.Substitute(tt, defaultMapping)
			assert.ErrorContains(t, err, "Invalid template")
		})
	}
}

// see https://github.com/docker/compose/issues/8601
func TestNonBraced(t *testing.T) {
	t.Parallel()

	substituted, warn, err := templatepkg.Substitute("$FOO-bar", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "first-bar", substituted)
}

func TestNoValueNoDefault(t *testing.T) {
	t.Parallel()

	{
		template := "This ${missing} var"
		result, warn, err := templatepkg.Substitute(template, defaultMapping)
		require.ErrorContains(t, warn, `The "missing" variable is not set. Defaulting to a blank string`)
		require.NoError(t, err)
		assert.Equal(t, "This  var", result)
	}

	{
		template := "This ${BAR} var"
		result, warn, err := templatepkg.Substitute(template, defaultMapping)
		require.NoError(t, warn)
		require.NoError(t, err)
		assert.Equal(t, "This  var", result)
	}
}

func TestValueNoDefault(t *testing.T) {
	t.Parallel()

	for _, template := range []string{"This $FOO var", "This ${FOO} var"} {
		result, warn, err := templatepkg.Substitute(template, defaultMapping)
		assert.NoError(t, warn)
		assert.NoError(t, err)
		assert.Equal(t, "This first var", result)
	}
}

func TestNoValueWithDefault(t *testing.T) {
	t.Parallel()

	for _, template := range []string{"ok ${missing:-def}", "ok ${missing-def}"} {
		result, warn, err := templatepkg.Substitute(template, defaultMapping)
		assert.NoError(t, warn)
		assert.NoError(t, err)
		assert.Equal(t, "ok def", result)
	}
}

func TestEmptyValueWithSoftDefault(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${BAR:-def}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok def", result)
}

func TestValueWithSoftDefault(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${FOO:-def}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok first", result)
}

func TestEmptyValueWithHardDefault(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${BAR-def}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok ", result)
}

func TestPresentValueWithUnset(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${UNSET_VAR:+presence_value}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok ", result)
}

func TestPresentValueWithUnset2(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${UNSET_VAR+presence_value}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok ", result)
}

func TestPresentValueWithNonEmpty(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${FOO:+presence_value}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok presence_value", result)
}

func TestPresentValueAndNonEmptyWithNonEmpty(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${FOO+presence_value}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok presence_value", result)
}

func TestPresentValueWithSet(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${BAR+presence_value}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok presence_value", result)
}

func TestPresentValueAndNotEmptyWithSet(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${BAR:+presence_value}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok ", result)
}

func TestNonAlphanumericDefault(t *testing.T) {
	t.Parallel()

	result, warn, err := templatepkg.Substitute("ok ${BAR:-/non:-alphanumeric}", defaultMapping)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok /non:-alphanumeric", result)
}

func TestInterpolationExternalInterference(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		template string
		expected string
	}{
		{
			template: "-ok ${BAR:-defaultValue}",
			expected: "-ok defaultValue",
		},
		{
			template: "+ok ${UNSET:-${BAR-defaultValue}}",
			expected: "+ok ",
		},
		{
			template: "-ok ${FOO:-defaultValue}",
			expected: "-ok first",
		},
		{
			template: ":-ok ${UNSET-defaultValue}",
			expected: ":-ok defaultValue",
		},
		{
			template: ":-ok ${BAR-defaultValue}",
			expected: ":-ok ",
		},
		{
			template: ":?ok ${BAR-defaultValue}",
			expected: ":?ok ",
		},
		{
			template: ":?ok ${BAR:-defaultValue}",
			expected: ":?ok defaultValue",
		},
		{
			template: ":+ok ${BAR:-defaultValue}",
			expected: ":+ok defaultValue",
		},
		{
			template: "+ok ${BAR-defaultValue}",
			expected: "+ok ",
		},
		{
			template: "?ok ${BAR:-defaultValue}",
			expected: "?ok defaultValue",
		},
	}

	for i, tt := range testCases {
		tt := tt

		t.Run(fmt.Sprintf("Interpolation Should not be impacted by outer text: %d", i), func(t *testing.T) {
			t.Parallel()

			result, warn, err := templatepkg.Substitute(tt.template, defaultMapping)
			assert.NoError(t, warn)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultsWithNestedExpansion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		template string
		expected string
	}{
		{
			template: "ok ${UNSET_VAR-$FOO}",
			expected: "ok first",
		},
		{
			template: "ok ${UNSET_VAR-${FOO}}",
			expected: "ok first",
		},
		{
			template: "ok ${UNSET_VAR-${FOO} ${FOO}}",
			expected: "ok first first",
		},
		{
			template: "ok ${BAR:-$FOO}",
			expected: "ok first",
		},
		{
			template: "ok ${BAR:-${FOO}}",
			expected: "ok first",
		},
		{
			template: "ok ${BAR:-${FOO} ${FOO}}",
			expected: "ok first first",
		},
		{
			template: "ok ${BAR+$FOO}",
			expected: "ok first",
		},
		{
			template: "ok ${BAR+$FOO ${FOO:+second}}",
			expected: "ok first second",
		},
	}

	for i, tt := range testCases {
		tt := tt

		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()

			result, warn, err := templatepkg.Substitute(tt.template, defaultMapping)
			assert.NoError(t, warn)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMandatoryVariableErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		template      string
		expectedError string
	}{
		{
			template:      "not ok ${UNSET_VAR:?Mandatory Variable Unset}",
			expectedError: "required variable UNSET_VAR is missing a value: Mandatory Variable Unset",
		},
		{
			template:      "not ok ${BAR:?Mandatory Variable Empty}",
			expectedError: "required variable BAR is missing a value: Mandatory Variable Empty",
		},
		{
			template:      "not ok ${UNSET_VAR:?}",
			expectedError: "required variable UNSET_VAR is missing a value",
		},
		{
			template:      "not ok ${UNSET_VAR?Mandatory Variable Unset}",
			expectedError: "required variable UNSET_VAR is missing a value: Mandatory Variable Unset",
		},
		{
			template:      "not ok ${UNSET_VAR?}",
			expectedError: "required variable UNSET_VAR is missing a value",
		},
	}

	for _, tc := range testCases {
		_, warn, err := templatepkg.Substitute(tc.template, defaultMapping)
		require.NoError(t, warn)
		require.ErrorContains(t, err, tc.expectedError)

		assert.ErrorAs(t, err, &MissingRequiredError)
	}
}

func TestMandatoryVariableErrorsWithNestedExpansion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		template      string
		expectedError string
	}{
		{
			template:      "not ok ${UNSET_VAR:?Mandatory Variable ${FOO}}",
			expectedError: "required variable UNSET_VAR is missing a value: Mandatory Variable first",
		},
		{
			template:      "not ok ${UNSET_VAR?Mandatory Variable ${FOO}}",
			expectedError: "required variable UNSET_VAR is missing a value: Mandatory Variable first",
		},
	}

	for _, tc := range testCases {
		_, _, err := templatepkg.Substitute(tc.template, defaultMapping)
		spew.Dump(err)
		require.ErrorContains(t, err, tc.expectedError)

		assert.ErrorAs(t, err, &MissingRequiredError)
	}
}

func TestDefaultsForMandatoryVariables(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		template string
		expected string
	}{
		{
			template: "ok ${FOO:?err}",
			expected: "ok first",
		},
		{
			template: "ok ${FOO?err}",
			expected: "ok first",
		},
		{
			template: "ok ${BAR?err}",
			expected: "ok ",
		},
	}

	for _, tc := range testCases {
		result, warn, err := templatepkg.Substitute(tc.template, defaultMapping)
		assert.NoError(t, warn)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, result)
	}
}

func TestSubstituteWithCustomFunc(t *testing.T) {
	t.Parallel()

	errIsMissing := func(substitution string, mapping templatepkg.Mapping) (string, bool, error) {
		value, found := mapping(substitution)
		if !found {
			return "", true, &templatepkg.InvalidTemplateError{
				Template: fmt.Sprintf("required variable %s is missing a value", substitution),
			}
		}

		return value, true, nil
	}

	result, warn, err := templatepkg.SubstituteWith("ok ${FOO}", defaultMapping, errIsMissing)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok first", result)

	result, warn, err = templatepkg.SubstituteWith("ok ${BAR}", defaultMapping, errIsMissing)
	assert.NoError(t, warn)
	assert.NoError(t, err)
	assert.Equal(t, "ok ", result)

	_, _, err = templatepkg.SubstituteWith("ok ${NOTHERE}", defaultMapping, errIsMissing)
	assert.ErrorContains(t, err, "required variable")
}

// TestPrecedence tests is the precedence on '-' and '?' is of the first match
func TestPrecedence(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		template string
		expected string
		err      error
	}{
		{
			template: "${UNSET_VAR?bar-baz}", // Unexistent variable
			expected: "",
			err: &templatepkg.MissingRequiredError{
				Variable: "UNSET_VAR",
				Reason:   "bar-baz",
			},
		},
		{
			template: "${UNSET_VAR-myerror?msg}", // Unexistent variable
			expected: "myerror?msg",
			err:      nil,
		},

		{
			template: "${FOO?bar-baz}", // Existent variable
			expected: "first",
		},
		{
			template: "${BAR:-default_value_for_empty_var}", // Existent empty variable
			expected: "default_value_for_empty_var",
		},
		{
			template: "${UNSET_VAR-default_value_for_unset_var}", // Unset variable
			expected: "default_value_for_unset_var",
		},
	}

	for _, tc := range testCases {
		result, warn, err := templatepkg.Substitute(tc.template, defaultMapping)
		require.NoError(t, warn)
		assert.Equal(t, tc.err, err)
		assert.Equal(t, tc.expected, result)
	}
}

func TestExtractVariables(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		dict     map[string]interface{}
		expected map[string]templatepkg.Variable
	}{
		{
			name:     "empty",
			dict:     map[string]interface{}{},
			expected: map[string]templatepkg.Variable{},
		},
		{
			name: "no-variables",
			dict: map[string]interface{}{
				"foo": "bar",
			},
			expected: map[string]templatepkg.Variable{},
		},
		{
			name: "variable-without-curly-braces",
			dict: map[string]interface{}{
				"foo": "$bar",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar"},
			},
		},
		{
			name: "variable",
			dict: map[string]interface{}{
				"foo": "${bar}",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar", DefaultValue: ""},
			},
		},
		{
			name: "required-variable",
			dict: map[string]interface{}{
				"foo": "${bar?:foo}",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar", DefaultValue: "", Required: true},
			},
		},
		{
			name: "required-variable2",
			dict: map[string]interface{}{
				"foo": "${bar?foo}",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar", DefaultValue: "", Required: true},
			},
		},
		{
			name: "default-variable",
			dict: map[string]interface{}{
				"foo": "${bar:-foo}",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar", DefaultValue: "foo"},
			},
		},
		{
			name: "default-variable2",
			dict: map[string]interface{}{
				"foo": "${bar-foo}",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar", DefaultValue: "foo"},
			},
		},
		{
			name: "multiple-values",
			dict: map[string]interface{}{
				"foo": "${bar:-foo}",
				"bar": map[string]interface{}{
					"foo": "${fruit:-banana}",
					"bar": "vegetable",
				},
				"baz": []interface{}{
					"foo",
					"$docker:${project:-cli}",
					"$toto",
				},
			},
			expected: map[string]templatepkg.Variable{
				"bar":     {Name: "bar", DefaultValue: "foo"},
				"fruit":   {Name: "fruit", DefaultValue: "banana"},
				"toto":    {Name: "toto", DefaultValue: ""},
				"docker":  {Name: "docker", DefaultValue: ""},
				"project": {Name: "project", DefaultValue: "cli"},
			},
		},
		{
			name: "presence-value-nonEmpty",
			dict: map[string]interface{}{
				"foo": "${bar:+foo}",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar", PresenceValue: "foo"},
			},
		},
		{
			name: "presence-value",
			dict: map[string]interface{}{
				"foo": "${bar+foo}",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar", PresenceValue: "foo"},
			},
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := templatepkg.ExtractVariables(tt.dict)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSubstitutionFunctionChoice(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name   string
		input  string
		symbol string
	}{
		{"Error when EMPTY or UNSET", "VARNAME:?val?ue", ":?"},
		{"Error when UNSET 1", "VARNAME?val:?ue", "?"},
		{"Error when UNSET 2", "VARNAME?va-lu+e:?e", "?"},
		{"Error when UNSET 3", "VARNAME?va+lu-e:?e", "?"},

		{"Default when EMPTY or UNSET", "VARNAME:-value", ":-"},
		{"Default when UNSET 1", "VARNAME-va:-lu:?e", "-"},
		{"Default when UNSET 2", "VARNAME-va+lu?e", "-"},
		{"Default when UNSET 3", "VARNAME-va?lu+e", "-"},

		{"Default when NOT EMPTY", "VARNAME:+va:?lu:-e", ":+"},
		{"Default when SET 1", "VARNAME+va:+lue", "+"},
		{"Default when SET 2", "VARNAME+va?lu-e", "+"},
		{"Default when SET 3", "VARNAME+va-lu?e", "+"},
	}

	for _, tt := range testcases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			symbol, _ := templatepkg.GetSubstitutionFunctionForTemplate(tt.input)
			assert.Equal(t, symbol, tt.symbol,
				fmt.Sprintf("Wrong on output for: %s got symbol -> %#v", tt.input, symbol),
			)
		})
	}
}

func TestNoValueWithCurlyBracesDefault(t *testing.T) {
	t.Parallel()

	for _, template := range []string{`ok ${missing:-{"json":1}}`, `ok ${missing-{"json":1}}`} {
		result, warn, err := templatepkg.Substitute(template, defaultMapping)
		assert.NoError(t, warn)
		assert.NoError(t, err)
		assert.Equal(t, `ok {"json":1}`, result)
	}
}

// TODO: figure out whats up with this one later
// func TestValueWithCurlyBracesDefault(t *testing.T) {
// 	t.Parallel()

// 	for _, template := range []string{`ok ${JSON:-{"json":1}}`, `ok ${JSON-{"json":1}}`} {
// 		result, warn, err := templatepkg.Substitute(template, defaultMapping)
// 		assert.NoError(t, warn)
// 		assert.NoError(t, err)
// 		assert.Equal(t, `ok {"json":2}`, result)
// 	}
// }
