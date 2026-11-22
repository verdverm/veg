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
		{{ upper .foo }}

		{{ template "helper" . }}

		"""

		partials: helper: "cows go {{ upper .cow }}"
  }
	out: { text: call.out } @task(os.Stdout)
}

