package upsert

// Setting is a bitmask for controlling Upsert behavior
type Setting int

const (
	// SkipIfSame will skip the upsert operation if the incoming KEY+VALUE is identical to the one in the document.
	// This is mostly used in the [update] command where it would not emit a "changed" event.
	SkipIfSame Setting = 1 << iota

	// SkipIfSet will skip the upsert operation if the KEY exists in the document.
	// This is useful for adding new KEY to the config file (e.g. during migration/upgrade), but never changing the key
	// if it already exists, regardless of the VALUE/Assignment configuration
	SkipIfSet

	// ErrorIfMissing will abort if the KEY does *NOT* exists in the target document. This is useful for ensuring
	// you only update keys, not accidentally creating new keys (e.g. in case of a typo)
	ErrorIfMissing
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
