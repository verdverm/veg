package agentic

import (
	"github.com/verdverm/veg/schemas"
)

Model: {
  schemas.#veg

  #veg: agentic: { root: true, kind: "model" }

  name: string
  id:   string

  features?: _
  prices?: _
}