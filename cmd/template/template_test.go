package template_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestTemplateCommand(t *testing.T) {
	test_helpers.RunFileBasedCommandTests(t, 0, "template")
}
