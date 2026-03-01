package template

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strings"

	slogctx "github.com/veqryn/slog-context"
	"mvdan.cc/sh/v3/expand"
)

// Resolver is a user-supplied function which maps from variable names to values.
// Returns the value as a string and a bool indicating whether
// the value is present, to distinguish between an empty string
// and the absence of a value.
type Resolver func(string) (string, bool)

type AccessibleVariables func() map[string]string

type EnvironmentHelper struct {
	Resolver            Resolver
	AccessibleVariables AccessibleVariables
	MissingKeyCallback  func(string)
}

func (helper EnvironmentHelper) Get(name string) expand.Variable {
	if val, ok := helper.Resolver(name); ok {
		return expand.Variable{
			Set:      true,
			Str:      val,
			Exported: true,
			ReadOnly: false,
			Kind:     expand.String,
		}
	}

	if val, ok := os.LookupEnv(name); ok {
		return expand.Variable{
			Set:      true,
			Str:      val,
			Exported: true,
			Kind:     expand.String,
		}
	}

	switch name {
	case "UID", "EUID":
		user, _ := user.Current()

		return expand.Variable{
			Set:      true,
			Str:      user.Uid,
			Exported: true,
			Kind:     expand.String,
		}

	case "GID":
		user, _ := user.Current()

		return expand.Variable{
			Set:      true,
			Str:      user.Gid,
			Exported: true,
			Kind:     expand.String,
		}

	case "IFS":
		return expand.Variable{
			Set:      true,
			Str:      `$' \t\n\C-@'`,
			Exported: true,
			Kind:     expand.String,
		}

	case "OPTIND":
		return expand.Variable{
			Set:      true,
			Str:      `1`,
			Exported: true,
			Kind:     expand.String,
		}
	}

	helper.MissingKeyCallback(name)

	return expand.Variable{
		Kind: expand.Unset,
	}
}

func (l EnvironmentHelper) Each(callback func(name string, vr expand.Variable) bool) {
	for k, v := range l.AccessibleVariables() {
		callback(k, expand.Variable{
			Set:      true,
			Str:      v,
			Exported: true,
			ReadOnly: false,
			Kind:     expand.String,
		})
	}

	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)

		callback(parts[0], expand.Variable{
			Set:      true,
			Str:      parts[1],
			Exported: true,
			ReadOnly: false,
			Kind:     expand.String,
		})
	}
}

func DefaultMissingKeyCallback(ctx context.Context, input string) func(string) {
	variables := ExtractVariables(ctx, input)

	return func(key string) {
		variable, ok := variables[key]

		// shouldn't be a lookup for anything that
		if !ok {
			slogctx.Warn(ctx, fmt.Sprintf("The [ $%s ] key is not set. Defaulting to a blank string.", key))

			return
		}

		// Required variables are errors, so we ignore them as warnings
		if variable.Required {
			return
		}

		// If the variable has a default value, then it's not missing
		if len(variable.DefaultValue) > 0 {
			return
		}

		// If the variable has a alternate/presence value, then it's not missing
		if len(variable.PresenceValue) > 0 {
			return
		}

		slogctx.Warn(ctx, fmt.Sprintf("The [ $%s ] key is not set. Defaulting to a blank string.", key))
	}
}
