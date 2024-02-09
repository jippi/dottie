package shared

import "github.com/spf13/pflag"

func BoolFlag(flags *pflag.FlagSet, name string) bool {
	val, _ := flags.GetBool(name)

	return val
}

func StringFlag(flags *pflag.FlagSet, name string) string {
	val, _ := flags.GetString(name)

	return val
}
