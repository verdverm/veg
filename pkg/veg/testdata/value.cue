package hack

val: {
	foo: "bar"

	$veg: {
		metadata: {
			labels: {
				app: "val-app"
				env: "dev"
			}
		}
		datamodel: {
			root:    true
			history: true
		}
	}
}
