package template_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestTemplateCommand(t *testing.T) {
	t.Parallel()

	test_helpers.RunFileBasedCommandTests(t, 0, "template")
}
