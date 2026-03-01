package print_cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldColorOutput(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		env      map[string]string
		expected bool
	}{
		{name: "default", args: nil, expected: false},
		{name: "color", args: []string{"--color"}, expected: true},
		{name: "no-color", args: []string{"--no-color"}, expected: false},
		{name: "color false", args: []string{"--color=false"}, expected: false},
		{name: "pretty", args: []string{"--pretty"}, expected: true},
		{name: "pretty color", args: []string{"--pretty", "--color"}, expected: true},
		{name: "pretty no-color", args: []string{"--pretty", "--no-color"}, expected: false},
		{name: "pretty color false", args: []string{"--pretty", "--color=false"}, expected: false},
		{name: "no color env + color", args: []string{"--color"}, env: map[string]string{"NO_COLOR": "1"}, expected: false},
		{name: "no color env + pretty", args: []string{"--pretty"}, env: map[string]string{"NO_COLOR": "1"}, expected: false},
		{name: "no color env + pretty color", args: []string{"--pretty", "--color"}, env: map[string]string{"NO_COLOR": "1"}, expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.env {
				t.Setenv(key, value)
			}

			command := New()
			require.NoError(t, command.ParseFlags(tc.args))
			assert.Equal(t, tc.expected, shouldColorOutput(command))
		})
	}
}
