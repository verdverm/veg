package hack

import (
	"github.com/verdverm/veg/schemas"
	"github.com/verdverm/veg/schemas/dm"
	"github.com/verdverm/veg/schemas/dm/fields"
)

// anything with `#veg: history: true`
// will have History and Lacuna enrichments
// these will be reflected in various parts
// of the datamodel and code gen

Datamodel: dm.Datamodel & {
	// implied through definition, duplicated here for example clarity
	$veg: metadata: {
		id:   "datamodel-abc123"
		name: "MyDatamodel"
	}
	// entire value / datamodel has history without any extra annotation
	// the config below would override that

	// permanent id and changable name
	// this makes renaming a table & type possible

	Config: {
		// track a full-object as a CUE value
		$veg: {
			history: true
			cue:     true
		}

		host: string
		port: string

		dbconn: string
	}

	// these are the models for the application
	// they can map onto database tables or apis
	Models: {
		// create point where veg can list, info, etc..
		#veg: datamodel: node: true
		@node()

		#veg: datamodel: ordered: true

		@ordered() // for stability, see below

		// each struct field is a Model
		[string]~(N, _): Model & {#veg: metadata: name: N}

		// Actual Models
		"User": User
	}
}

Model: M={
	schemas.#veg// needed for reFerences

	#veg: datamodel: history: true // needed for CUE compat

	@history() // veg only shorthand

	// Lineage fields will be filled by veg
	// #veg: Lense: ...
	// #veg: History: ...

	// for easy access
	Name: M.#veg.metadata.name

	// These are the fields of a model
	// they can map onto database columnts and form fields
	Fields: {
		// create point where veg can list, info, etc..
		#veg: datamodel: node: true
		@node()

		// for stability, see below
		#veg: datamodel: ordered: true

		@ordered() // shorthand

		[string]~(N, _): Field & {#veg: metadata: name: N}
	}

	// OrderedFields: [...Fields] will be
	// inject here for order stability

	// if we want Relations as a separate value
	// we can process the fields to extract them
}

Field: F={
	schemas.#veg// needed for reFerences
	#veg: datamodel: history: true

	@history() // shorthand

	// Lineage fields will be filled by veg
	// #veg: Lense: ...
	// #veg: History: ...

	// for easy access
	Name: F.#veg.metadata.name
	Type: string

	// relation type, open to be flexible
	Reln?: string

	// what about {val, *val, []val, []*val}
	// we probably don't care about pointer here
	//   that is a language detail (code gen target)

	// we can enrich this for various types
	// in our app or other reusable datamodels
}

User: {
	@id(437df0a6-a57c-11ed-82e6-00163ece1931)
	Fields: {
		ID:        fields.UUID
		CreatedAt: fields.Datetime
		UpdatedAt: fields.Datetime
		DeletedAt: fields.Datetime

		email:    fields.Email
		password: fields.Password
		active:   fields.Bool

		// this is the new field
		username: fields.String

		persona: fields.Enum & {
			Vals: ["guest", "user", "admin", "owner"]
			Default: "guest"
		}
	}
}
