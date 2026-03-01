package print_cmd_test

import (
	"testing"

	print_cmd "github.com/jippi/dottie/cmd/print"
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

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for key, value := range testCase.env {
				t.Setenv(key, value)
			}

			command := print_cmd.New()
			require.NoError(t, command.ParseFlags(testCase.args))
			assert.Equal(t, testCase.expected, print_cmd.ShouldColorOutput(command))
		})
	}
}
