package shared

import (
	"os"

	"github.com/spf13/pflag"
)

func ColorEnabled(flags *pflag.FlagSet, name string) bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}

	return BoolWithInverseValue(flags, name)
}
