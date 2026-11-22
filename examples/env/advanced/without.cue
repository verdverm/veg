package advanced

import (
	"github.com/verdverm/veg/catalogs/env/bases"
	"github.com/verdverm/veg/schemas/env"
)

without: env.#Container & {
  @env(without)
  from: bases.rocky8.minimal
  steps: [
    env.WithoutFile & { path: "/usr/bin/curl" }
  ]
}