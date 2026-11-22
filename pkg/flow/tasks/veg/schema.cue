package veg

Tempate: {
	@task(veg.Template)

	name: string | *""
	data: _

	template: string
	partials: [string]: string

	out: string
}
