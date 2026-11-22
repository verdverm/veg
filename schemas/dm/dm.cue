package dm

import (
	"github.com/verdverm/veg/schemas"
)

// This is a complete Value tracked as one
// useful for schemas, config, and NoSQL
Object: {
	schemas.#veg // needed for reFerences
	#veg: datamodel: root: true

	TrackHistory

	// all fields will be tracked
}

// This is like object, but supports cue values
// (todo, should support full lattice)
Value: {
	Object
	#veg: datamodel: cue: true
}

// This is a general datamodel useful in many applications
// It can be expanded and enriched to cover more
// Useful for SQL, APIs, forms, and similar
Datamodel: {
	schemas.#veg // needed for references
	#veg: datamodel: root: true
}

// Schema for a snapshot, can include anything else
Snapshot: {
	Timestamp: string | *""
}

// convenience type for embedding history
History: [...Snapshot]

TrackHistory: {
	#veg: datamodel: history: true // needed for CUE compat
	"Snapshot": Snapshot
	"History":  History
}
