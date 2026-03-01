package test_helpers_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestSettingHas(t *testing.T) {
	t.Parallel()

	if test_helpers.Setting(0).Has(test_helpers.ReadOnly) {
		t.Fatal("expected empty bitmask to not contain ReadOnly")
	}

	if !test_helpers.ReadOnly.Has(test_helpers.ReadOnly) {
		t.Fatal("expected ReadOnly bitmask to contain ReadOnly")
	}
}
