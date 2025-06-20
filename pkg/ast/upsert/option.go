package upsert

import (
	"errors"
	"fmt"
)

// Option is used to configure the [Upserter]
type Option func(*Upserter) error

// WithGroup configures the [Upserter] to the new [ast.Assignment] to this
// group when creating it within the [ast.Document]
func WithGroup(name string) Option {
	return func(upserter *Upserter) error {
		upserter.group = name

		return nil
	}
}

// EnableSettingIf will, depending on [boolean], either enable or disable a
// [Upserter] setting. Its mainly a convenience function to avoid if/else on the caller
// side - such as in cases of bool CLI flags changing controlling a setting.
func EnableSettingIf(setting Setting, boolean bool) Option {
	if boolean {
		return EnableSetting(setting)
	}

	return DisableSetting(setting)
}

// EnableSetting will set the provided [Setting] in the [Upserter] settings bitmask.
func EnableSetting(setting Setting) Option {
	return func(upserter *Upserter) error {
		upserter.settings = upserter.settings | setting

		return nil
	}
}

// WithSetting will remove the provided [Setting] in the [Upserter] settings bitmask.
func DisableSetting(setting Setting) Option {
	return func(upserter *Upserter) error {
		upserter.settings = upserter.settings &^ setting

		return nil
	}
}

// WithPlacementRelativeToKey configures the [Upserter] to add a new KEY in a
// specific place (in relation to the provided KEY) within a document, if the KEY does not already exists.
func WithPlacementRelativeToKey(placement Placement, key string) Option {
	return func(upserter *Upserter) error {
		if len(key) == 0 {
			return errors.New("empty 'KEY' was provided to placement logic")
		}

		if !placement.RequiresKey() {
			return fmt.Errorf("the placement (%s) does not support a relative KEY (%s), please use [WithPlacement] instead", placement, key)
		}

		other := upserter.document.Get(key)
		if other == nil {
			return fmt.Errorf("the KEY [%s] does not exists in the document", key)
		}

		upserter.placement = placement
		upserter.placementValue = key

		return nil
	}
}

// WithPlacementKey configures the [Upserter] to add a new KEY in a relative place within
// a document, if the KEY does not already exists.
func WithPlacement(placement Placement) Option {
	return func(upserter *Upserter) error {
		if placement.RequiresKey() {
			return fmt.Errorf("the Placement (%s) does requires a KEY, please use [WithPlacementKey] instead", placement)
		}

		upserter.placement = placement

		return nil
	}
}

// WithPlacementInGroup is like [WithPlacement] but further more configures the
// KEY's group (if any) as well
func WithPlacementInGroup(placement Placement, key string) Option {
	return func(upserter *Upserter) error {
		// First configure the placement
		if err := WithPlacementRelativeToKey(placement, key)(upserter); err != nil {
			return err
		}

		// Then force the group placement (if any)
		other := upserter.document.Get(key)
		if other.Group != nil {
			return WithGroup(other.Group.String())(upserter)
		}

		return nil
	}
}

// WithPlacementIgnoringEmpty is like [WithPlacement] but is NOOP if
// the KEY is an empty string.
// Mostly useful convenience method for passing through CLI flags directly.
func WithPlacementIgnoringEmpty(placement Placement, key string) Option {
	if len(key) == 0 {
		return func(u *Upserter) error {
			return nil
		}
	}

	return WithPlacementRelativeToKey(placement, key)
}

// WithPlacementInGroupIgnoringEmpty is like [WithPlacementInGroup] but is NOOP if
// the KEY is an empty string.
// Mostly useful convenience method for passing through CLI flags directly.
func WithPlacementInGroupIgnoringEmpty(placement Placement, key string) Option {
	if len(key) == 0 {
		return func(u *Upserter) error {
			return nil
		}
	}

	return WithPlacementInGroup(placement, key)
}

// WithSkipValidationRule allow you to skip/ignore specific validation rules
func WithSkipValidationRule(rules ...string) Option {
	return func(u *Upserter) error {
		u.ignoreValidationRules = rules

		return nil
	}
}
