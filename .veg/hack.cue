package config

import (
	// "github.com/verdverm/veg/catalogs/env/bases"
	"github.com/verdverm/veg/catalogs/env/packs"
	"github.com/verdverm/veg/schemas/env"
)

hack: {
	// used as a simple reproducer to determine that veg-dagger-engine
	// was missing a volume mount and using crazy amounts of disk
	diskUsage: env.#Container & {
		@env(hack-diskUsage)
		from: "ghcr.io/verdverm/veg-veg:v0.7.0-alpha.1"
		steps: [
			env.Mount & {path: "/work", source: src.code},
		]
	}

	sbom: {
		cuefig: env.#CuefigSBOM & {
			@env(hack-cuefig-sbom)
			path:   "hack-cuefig-sbom.cue"
			format: "cue"
			data:   diskUsage
		}
	}

	install: {
		crane: env.#ExportDir & {
			@env(hack-install-crane)
			path: "/usr/local/bin"
			sources: [(packs.containers.crane & {#distro: "Darwin"}).files]
			// sources: [files]
			include: ["crane", "krane"]
			wipe: true
		}
	}
}
