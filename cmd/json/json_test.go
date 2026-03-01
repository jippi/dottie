package json_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestJsonCommand(t *testing.T) {
	test_helpers.RunFileBasedCommandTests(t, 0, "json")
}
