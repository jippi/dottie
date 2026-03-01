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
	"errors"
	"fmt"
	"testing"

	templatepkg "github.com/jippi/dottie/pkg/template"
	"github.com/jippi/dottie/pkg/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaults = map[string]string{
	"FOO":  "first",
	"BAR":  "",
	"JSON": `{"json":2}`,
}

func defaultMapping(name string) (string, bool) {
	val, ok := defaults[name]

	return val, ok
}

func accessibleVariables() map[string]string {
	return defaults
}

// func TestEscaped(t *testing.T) {
// 	t.Parallel()
//
// 	result,  err := templatepkg.Substitute("$${foo}", defaultMapping)
// 	require.NoError(t, warn)
// 	require.NoError(t, err)
// 	assert.Equal(t, "${foo}", result)
// }

func TestSubstituteNoMatch(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "foo", defaultMapping, accessibleVariables)
	require.NoError(t, err)
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
		actual, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), expected, defaultMapping, accessibleVariables)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestInvalid(t *testing.T) {
	t.Parallel()

	invalidTemplates := []string{
		"${",
		"${}",
		"${ }",
		"${ foo}",
		// "${foo }",
		"${foo!}",
	}

	for i, tt := range invalidTemplates {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()

			_, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tt, defaultMapping, accessibleVariables)

			assert.ErrorContains(t, err, "Invalid template")
		})
	}
}

func TestSubstitute_MalformedArithmeticExpansionReturnsError(t *testing.T) {
	t.Parallel()

	malformed := "$(($'\\\"0\\\"0\\\"0\\\"0\\00'))"

	assert.NotPanics(t, func() {
		_, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), malformed, defaultMapping, accessibleVariables)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Invalid template")
	})
}

// see https://github.com/docker/compose/issues/8601
func TestNonBraced(t *testing.T) {
	t.Parallel()

	substituted, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "$FOO-bar", defaultMapping, accessibleVariables)
	require.NoError(t, err)
	assert.Equal(t, "first-bar", substituted)
}

func TestNoValueNoDefault(t *testing.T) {
	t.Parallel()

	{
		template := "This ${missing} var"
		result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), template, defaultMapping, accessibleVariables)

		// require.ErrorContains(t, warn, `The [ $missing ] key is not set. Defaulting to a blank string.`)
		require.NoError(t, err)
		assert.Equal(t, "This  var", result)
	}

	{
		template := "This ${BAR} var"
		result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), template, defaultMapping, accessibleVariables)

		require.NoError(t, err)
		assert.Equal(t, "This  var", result)
	}
}

func TestValueNoDefault(t *testing.T) {
	t.Parallel()

	for _, template := range []string{"This $FOO var", "This ${FOO} var"} {
		result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), template, defaultMapping, accessibleVariables)
		require.NoError(t, err)
		assert.Equal(t, "This first var", result)
	}
}

func TestNoValueWithDefault(t *testing.T) {
	t.Parallel()

	for _, template := range []string{"ok ${missing:-def}", "ok ${missing-def}"} {
		result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), template, defaultMapping, accessibleVariables)
		require.NoError(t, err)
		assert.Equal(t, "ok def", result)
	}
}

func TestEmptyValueWithSoftDefault(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${BAR:-def}", defaultMapping, accessibleVariables)
	require.NoError(t, err)
	assert.Equal(t, "ok def", result)
}

func TestValueWithSoftDefault(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${FOO:-def}", defaultMapping, accessibleVariables)
	require.NoError(t, err)
	assert.Equal(t, "ok first", result)
}

func TestEmptyValueWithHardDefault(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${BAR-def}", defaultMapping, accessibleVariables)
	require.NoError(t, err)
	assert.Equal(t, "ok ", result)
}

func TestPresentValueWithUnset(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${UNSET_VAR:+presence_value}", defaultMapping, accessibleVariables)
	require.NoError(t, err, "error")
	assert.Equal(t, "ok ", result)
}

func TestPresentValueWithUnset2(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${UNSET_VAR+presence_value}", defaultMapping, accessibleVariables)
	require.NoError(t, err, "error")
	assert.Equal(t, "ok ", result)
}

func TestPresentValueWithNonEmpty(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${FOO:+presence_value}", defaultMapping, accessibleVariables)
	require.NoError(t, err)
	assert.Equal(t, "ok presence_value", result)
}

func TestPresentValueAndNonEmptyWithNonEmpty(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${FOO+presence_value}", defaultMapping, accessibleVariables)
	require.NoError(t, err)
	assert.Equal(t, "ok presence_value", result)
}

func TestPresentValueWithSet(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${BAR+presence_value}", defaultMapping, accessibleVariables)
	require.NoError(t, err)
	assert.Equal(t, "ok presence_value", result)
}

func TestPresentValueAndNotEmptyWithSet(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${BAR:+presence_value}", defaultMapping, accessibleVariables)
	require.NoError(t, err)
	assert.Equal(t, "ok ", result)
}

func TestNonAlphanumericDefault(t *testing.T) {
	t.Parallel()

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${BAR:-/non:-alphanumeric}", defaultMapping, accessibleVariables)
	require.NoError(t, err)
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
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()

			result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tt.template, defaultMapping, accessibleVariables)
			require.NoError(t, err)
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
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()

			result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tt.template, defaultMapping, accessibleVariables)
			require.NoError(t, err)
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

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()

			_, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tt.template, defaultMapping, accessibleVariables)
			require.ErrorContains(t, err, tt.expectedError)

			missingRequiredError := &templatepkg.MissingRequiredError{}
			assert.ErrorAs(t, err, &missingRequiredError)
		})
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

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()

			_, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tt.template, defaultMapping, accessibleVariables)
			require.ErrorContains(t, err, tt.expectedError)

			missingRequiredError := &templatepkg.MissingRequiredError{}

			assert.ErrorAs(t, err, &missingRequiredError)
		})
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
		result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tc.template, defaultMapping, accessibleVariables)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, result)
	}
}

func TestSetVariableSemanticsForEmptyValue(t *testing.T) {
	t.Parallel()

	resolver := func(name string) (string, bool) {
		if name == "VAR" {
			return "", true
		}

		return "", false
	}

	accessible := func() map[string]string {
		return map[string]string{"VAR": ""}
	}

	result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${VAR+present}", resolver, accessible)
	require.NoError(t, err)
	assert.Equal(t, "ok present", result)

	result, err = templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), "ok ${VAR?err}", resolver, accessible)
	require.NoError(t, err)
	assert.Equal(t, "ok ", result)
}

// TestShellOperatorComplianceMatrix documents and verifies shell parameter
// expansion behavior for unset, empty, and non-empty variables.
//
// Operators covered:
//   - ${VAR-word}: use word only when VAR is unset
//   - ${VAR:-word}: use word when VAR is unset or empty
//   - ${VAR+word}: use word when VAR is set (including empty)
//   - ${VAR:+word}: use word when VAR is set and non-empty
//   - ${VAR?word}: error when VAR is unset
//   - ${VAR:?word}: error when VAR is unset or empty
func TestShellOperatorComplianceMatrix(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		vars      map[string]string
		template  string
		expected  string
		errSubstr string
	}

	testCases := []testCase{
		{name: "unset-dash", vars: map[string]string{}, template: "${VAR-default}", expected: "default"},
		{name: "unset-colon-dash", vars: map[string]string{}, template: "${VAR:-default}", expected: "default"},
		{name: "unset-plus", vars: map[string]string{}, template: "${VAR+alt}", expected: ""},
		{name: "unset-colon-plus", vars: map[string]string{}, template: "${VAR:+alt}", expected: ""},
		{name: "unset-question", vars: map[string]string{}, template: "${VAR?err}", errSubstr: "required variable VAR is missing a value: err"},
		{name: "unset-colon-question", vars: map[string]string{}, template: "${VAR:?err}", errSubstr: "required variable VAR is missing a value: err"},

		{name: "empty-dash", vars: map[string]string{"VAR": ""}, template: "${VAR-default}", expected: ""},
		{name: "empty-colon-dash", vars: map[string]string{"VAR": ""}, template: "${VAR:-default}", expected: "default"},
		{name: "empty-plus", vars: map[string]string{"VAR": ""}, template: "${VAR+alt}", expected: "alt"},
		{name: "empty-colon-plus", vars: map[string]string{"VAR": ""}, template: "${VAR:+alt}", expected: ""},
		{name: "empty-question", vars: map[string]string{"VAR": ""}, template: "${VAR?err}", expected: ""},
		{name: "empty-colon-question", vars: map[string]string{"VAR": ""}, template: "${VAR:?err}", errSubstr: "required variable VAR is missing a value: err"},

		{name: "value-dash", vars: map[string]string{"VAR": "value"}, template: "${VAR-default}", expected: "value"},
		{name: "value-colon-dash", vars: map[string]string{"VAR": "value"}, template: "${VAR:-default}", expected: "value"},
		{name: "value-plus", vars: map[string]string{"VAR": "value"}, template: "${VAR+alt}", expected: "alt"},
		{name: "value-colon-plus", vars: map[string]string{"VAR": "value"}, template: "${VAR:+alt}", expected: "alt"},
		{name: "value-question", vars: map[string]string{"VAR": "value"}, template: "${VAR?err}", expected: "value"},
		{name: "value-colon-question", vars: map[string]string{"VAR": "value"}, template: "${VAR:?err}", expected: "value"},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			resolver := func(name string) (string, bool) {
				value, ok := testCase.vars[name]

				return value, ok
			}

			accessible := func() map[string]string {
				return testCase.vars
			}

			result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), testCase.template, resolver, accessible)

			if testCase.errSubstr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.errSubstr)
				assert.Equal(t, "", result)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestComplexNestedInterpolationBehavior(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		vars     map[string]string
		template string
		expected string
	}

	testCases := []testCase{
		{
			name:     "nested-default-uses-next-fallback",
			vars:     map[string]string{},
			template: "${A:-${B:-fallback}}",
			expected: "fallback",
		},
		{
			name:     "nested-default-uses-inner-variable",
			vars:     map[string]string{"B": "inner"},
			template: "${A:-${B:-fallback}}",
			expected: "inner",
		},
		{
			name:     "nested-default-prefers-outer-variable",
			vars:     map[string]string{"A": "outer", "B": "inner"},
			template: "${A:-${B:-fallback}}",
			expected: "outer",
		},
		{
			name:     "nested-alternate-with-inner-default",
			vars:     map[string]string{"A": "set"},
			template: "${A:+${B:-alt}}",
			expected: "alt",
		},
		{
			name:     "nested-alternate-with-inner-value",
			vars:     map[string]string{"A": "set", "B": "value"},
			template: "${A:+${B:-alt}}",
			expected: "value",
		},
		{
			name:     "mixed-concatenation-with-defaults-and-alternate",
			vars:     map[string]string{"A": "left", "B": ""},
			template: "pre-${A}-${B:-middle}-${A:+tail}",
			expected: "pre-left-middle-tail",
		},
		{
			name:     "double-nested-operators",
			vars:     map[string]string{"A": "", "B": "bee", "C": "see"},
			template: "${A:-${B:+${C:-fallback}}}",
			expected: "see",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resolver := func(name string) (string, bool) {
				value, ok := tc.vars[name]

				return value, ok
			}

			accessible := func() map[string]string {
				return tc.vars
			}

			result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tc.template, resolver, accessible)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestComplexInterpolationSpecialConstructsRemainLiteral(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		template string
		expected string
	}

	testCases := []testCase{
		{
			name:     "default-containing-process-substitution-is-not-executed",
			template: "${UNSET:-<(cat /dev/null)}",
			expected: "<(cat /dev/null)",
		},
		{
			name:     "process-substitution-and-regular-interpolation",
			template: "prefix ${FOO} ${UNSET:-<(cat /dev/null)} suffix",
			expected: "prefix first <(cat /dev/null) suffix",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tc.template, defaultMapping, accessibleVariables)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
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
			template: "${UNSET_VAR?bar-baz}", // Nonexistent variable
			expected: "",
			err: &templatepkg.MissingRequiredError{
				Variable: "UNSET_VAR",
				Reason:   "bar-baz",
			},
		},
		{
			template: "${UNSET_VAR-myerror?msg}", // Nonexistent variable
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

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()

			result, err := templatepkg.Substitute(test_helpers.CreateTestContext(t, nil, nil), tt.template, defaultMapping, accessibleVariables)

			var asErr *templatepkg.MissingRequiredError

			switch {
			case errors.As(err, &asErr):
				assert.Equal(t, tt.err, asErr)

			default:
				assert.Equal(t, tt.err, err)
			}

			assert.Equal(t, tt.expected, result)
		})
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
		{
			name: "default-with-process-substitution",
			dict: map[string]interface{}{
				"foo": "${bar:-<(cat /dev/null)}",
			},
			expected: map[string]templatepkg.Variable{
				"bar": {Name: "bar", DefaultValue: "<(cat /dev/null)"},
			},
		},
		{
			name: "alternate-with-empty-value",
			dict: map[string]interface{}{
				"foo": "${0+}",
			},
			expected: map[string]templatepkg.Variable{
				"0": {Name: "0"},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := templatepkg.ExtractVariables(t.Context(), tt.dict)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
