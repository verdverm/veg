package sql

import (
	"github.com/verdverm/veg/schemas"
	"github.com/verdverm/veg/schemas/dm"
)

Datamodel: {
	#veg: datamodel: root: true

	schemas.#veg
	dm.TrackHistory

	// these are the models for the application
	// they can map onto database tables and apis
	Models: {
		#veg: datamodel: node:    true
		#veg: datamodel: ordered: true
		[string]~(N, _): Model & {Name: N, #veg: metadata: name: N}
	}

	// OrderedModels: [...Model] will be
	// inject here for order stability
}

Model: M={
	schemas.#veg
	dm.TrackHistory

	// for easy access
	Name:   M.#veg.metadata.name
	Plural: string | *"\(Name)s"

	// These are the fields of a model
	// they can map onto database columnts and form fields
	Fields: {
		#veg: datamodel: node:    true
		#veg: datamodel: ordered: true
		[string]~(N, _): Field & {Name: N, #veg: metadata: name: N}
	}

	// OrderedFields: [...Fields] will be
	// inject here for order stability

	// if we want Relations as a separate value
	// we can process the fields to extract them
}

Field: {
	// schemas.#veg
	// TODO, decide if these should be the default
	// #veg: datamodel: history: true // needed for CUE compat
	// History: dm.History

	Name: string
	Type: string

	// relation type, open to be flexible
	Relation?: {
		"Name": string | *Name
		Type:   "has-one" | "has-many" | "belongs-to" | "many-to-many"
		Other:  string // technically a cue path, but as a string
	}

	// what about {val, *val, []val, []*val}
	// we probably don't care about pointer here
	//   that is a language detail (code gen target)

	// we can enrich this for various types
	// in our app or other reusable datamodels
}
