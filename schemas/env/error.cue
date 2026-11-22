package env

import (
	"github.com/verdverm/veg/schemas"
)

#Error: {
	schemas.Hof
	#veg: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "error"
	}
	$kind: "#error"

  message?: string
}