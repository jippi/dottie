package upsert

import (
	"fmt"
	"strings"
)

// Setting is a bitmask for controlling Upsert behavior
type Setting int

const (
	// SkipIfSame will skip the upsert operation if the incoming KEY+VALUE is identical to the one in the document.
	// This is mostly used in the [update] command where it would not emit a "changed" event.
	SkipIfSame Setting = 1 << iota

	// SkipIfExists will skip the upsert operation if the KEY exists in the document (empty or not).
	// This is useful for adding new KEY to the config file (e.g. during migration/upgrade), but never changing the key
	// if it already exists, regardless of the VALUE/Assignment configuration
	SkipIfExists

	// SkipIfSet will skip the upsert operation if the KEY exists in the document and *NOT* empty.
	SkipIfSet

	// Validate the KEY/VALUE pair and fail the operation if its invalid
	Validate

	// ErrorIfMissing will abort if the KEY does *NOT* exists in the target document. This is useful for ensuring
	// you only update keys, not accidentally creating new keys (e.g. in case of a typo).
	ErrorIfMissing

	// Replace comments on *existing* Assignments.
	//
	// Normally comments are only applied to *NEW* keys and not on existing ones.
	UpdateComments

	// Only here to make iteration over the Settings list easier, is not used externally and have no special meaning.
	// See: settings.name()
	maxKey
)

// Has checks if [check] exists in the [settings] bitmask or not.
func (bitmask Setting) Has(setting Setting) bool {
	// If [settings] is 0, its an initialized/unconfigured bitmask, so no settings exists.
	//
	// This is true since all UpsertSetting starts from "1", not "0".
	if bitmask == 0 {
		return false
	}

	return bitmask&setting != 0
}

// Pretty-print Setting value, both for single and multiple keys
func (setting Setting) String() string {
	// Single key bitmask
	switch setting {
	case SkipIfSame, SkipIfExists, SkipIfSet, Validate, ErrorIfMissing, UpdateComments:
		return fmt.Sprintf("upsert.Setting<%s>", setting.name())
	case maxKey:
	}

	// Multi key bitmask
	var names []string

	for key := SkipIfSame; key < maxKey; key <<= 1 {
		if setting&key != 0 {
			names = append(names, key.name())
		}
	}

	return fmt.Sprintf("upsert.Setting<%s>", strings.Join(names, " | "))
}

// name returns the human readable name for an individual Setting key.
func (setting Setting) name() string {
	switch setting {
	case SkipIfSame:
		return "SkipIfSame"

	case SkipIfExists:
		return "SkipIfExists"

	case SkipIfSet:
		return "SkipIfSet"

	case Validate:
		return "Validate"

	case ErrorIfMissing:
		return "ErrorIfMissing"

	case UpdateComments:
		return "ReplaceComments"

	case maxKey:
	}

	return "UNKNOWN"
}
