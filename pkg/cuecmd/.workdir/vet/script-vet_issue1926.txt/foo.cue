package foo

_env: =~"(prd)-"

#Foo

#Foo: {
	env: "dev" | _env
	baz: bool | *false

	if env == "dev" {
		_defaults.dev
	}
}

_defaults: {
	"dev": {
		baz: true
	}
}

