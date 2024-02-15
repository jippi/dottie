package shared

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func BoolFlag(flags *pflag.FlagSet, name string) bool {
	val, _ := flags.GetBool(name)

	return val
}

func StringFlag(flags *pflag.FlagSet, name string) string {
	val, _ := flags.GetString(name)

	return val
}

func StringSliceFlag(flags *pflag.FlagSet, name string) []string {
	val, _ := flags.GetStringSlice(name)

	return val
}

func BoolWithInverse(cmd *cobra.Command, name string, value bool, usage, negativeUsage string) {
	cmd.Flags().Bool(name, value, usage)
	cmd.Flags().Bool("no-"+name, !value, negativeUsage)

	cmd.MarkFlagsMutuallyExclusive(name, "no-"+name)
}

func BoolWithInverseValue(flags *pflag.FlagSet, name string) bool {
	switch {
	// If "no-" flag was used, return that value
	case flags.Lookup("no-" + name).Changed:
		val, _ := flags.GetBool("no-" + name)

		return !val

		// Otherwise, use the default (positive) flag
	default:
		val, _ := flags.GetBool(name)

		return val
	}
}

// optional interface to indicate boolean flags that can be
// supplied without "=value" text
type boolFlag interface {
	pflag.Value
	IsBoolFlag() bool
}

// -- bool Value
type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val

	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)

	return err
}

func (b *boolValue) Type() string {
	return "bool"
}

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }
