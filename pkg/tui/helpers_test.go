package tui_test

import (
	"strings"
	"testing"

	"github.com/jippi/dottie/pkg/tui"
)

func TestTransformColorShadeTintAndFallback(t *testing.T) {
	t.Parallel()

	base := "#336699"
	shade := tui.TransformColor(base, "shade", 0.5)
	tint := tui.TransformColor(base, "tint", 0.5)
	fallback := tui.TransformColor(base, "unknown", 0.5)

	if shade == base || tint == base {
		t.Fatal("expected shade and tint to transform the base color")
	}

	if fallback != base {
		t.Fatalf("expected unknown filter to return base color, got %q", fallback)
	}

	if !strings.HasPrefix(shade, "#") || !strings.HasPrefix(tint, "#") {
		t.Fatal("expected transformed colors to be hex formatted")
	}
}

func TestTransformColorMixPanics(t *testing.T) {
	t.Parallel()

	didPanic := false

	defer func() {
		if recover() != nil {
			didPanic = true
		}

		if !didPanic {
			t.Fatal("expected mix filter to panic")
		}
	}()

	tui.TransformColor("#336699", "mix", 0.5)
}
