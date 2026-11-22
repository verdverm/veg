package advanced

import (
	"github.com/verdverm/veg/catalogs/env/bases"
	// "github.com/verdverm/veg/catalogs/env/utils"
	"github.com/verdverm/veg/schemas/env"
)

scratch: {
  prep: env.#Container & {
    @env(scratch-prep)
    from: bases.alpine3.minimal
  }
  image: env.#Container & {
    @env(scratch-image)
    from: "scratch"
  }
}