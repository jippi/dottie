package upsert

import (
	"fmt"
	"slices"
)

// Placement is a setting for deciding where a *new* KEY should be added within the document/group it's targeting
// When a KEY already exists, its placement will not be updated
type Placement uint

const (
	// NOTE: [AddLast] is the *default* when not configured in the [Upserter] (since its the 'empty' type of int aka '0')
	AddLast Placement = iota
	AddAfterKey
	AddBeforeKey
	AddFirst
)

// placementRequiresKey is the list of Placements that requires a valid KEY
// to satisfy its placement strategy.
var placementRequiresKey = []Placement{
	AddAfterKey,
	AddBeforeKey,
}

func (p Placement) RequiresKey() bool {
	return slices.Contains(placementRequiresKey, p)
}

func (p Placement) String() string {
	switch p {
	case AddLast:
		return "Placement<AddLast>"

	case AddFirst:
		return "Placement<AddFirst>"

	case AddAfterKey:
		return "Placement<AddAfterKey>"

	case AddBeforeKey:
		return "Placement<AddBeforeKey>"

	default:
		panic(fmt.Errorf("unexpected upsert.Placement value: %d", p))
	}
}
