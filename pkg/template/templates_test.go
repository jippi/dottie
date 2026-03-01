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
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := templatepkg.ExtractVariables(t.Context(), tt.dict)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
