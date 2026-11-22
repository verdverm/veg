package agentic

import (
	"github.com/verdverm/veg/schemas"
)

Tool: {
  schemas.#veg

  #veg: agentic: { root: true, kind: "tool" }

  name: string

}