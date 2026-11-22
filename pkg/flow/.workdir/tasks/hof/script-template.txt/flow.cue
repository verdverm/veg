package veg

tasks: {
  @flow()
	call: { 
    @task(veg.Template)
		data: {
			foo: "bar"
			cow: "moo"
		}
		template: """
		{{ .foo }}

		{{ template "helper" . }}

		"""

		partials: helper: "cows go {{ .cow }}"
  }
	out: { text: call.out } @task(os.Stdout)
}

