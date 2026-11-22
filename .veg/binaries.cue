package config

import (
	"github.com/verdverm/veg/catalogs/env/packs/lang"
	"github.com/verdverm/veg/schemas/env"
)

// go multi-target builds
// from source build too for local setup

ctr: {
	builder: env.#Container & {
		@env(veg-cli-builder)
		from: lang.go.ctr.base
		steps: [
			env.Dir & {path: "/adk", source: src.adk.fork},
			env.Dir & {path: "/dagger", source: src.dagger.fork},
			env.Dir & {path: "/work", source: env.#HostDir & {
				path: flags.disk
				name: "veg-bin-src"
				include: [
					"go.mod", "go.sum", "cue.mod",
					"cmd", "lib", "flow", "script",
				]
			}},
		]
	}
}

// File version
veg: File: {for k, v in veg.cli.matrix {
	(k): env.File & {
		#veg: {id: v.#veg.id, metadata: v.#veg.metadata}
		name:    v.name
		path:    string | *"/usr/local/bin/veg"
		content: v
	}
}}
// #File verions
veg: cli: {
	_maker: env.#File & {
		// metadata
		@env()
		#veg: id: "veg-cli-\(#variant)"
		#veg: metadata: name: #veg.id

		// params
		#variant: string
		#goos:    string
		#arch:    string

		path: "./bins/veg-\(#variant)"
		source: env.#Container & {
			from: ctr.builder
			steps: [
				env.EnvVars & {GOOS: #goos, GOARCH: #arch},
				env.Exec & {args: ["go", "build", "-ldflags", "-w", "-o", "./bins/veg-\(#variant)", "./cmd/veg"]},
			]
		}
	}

	local: _maker & {#variant: "local", #goos: flags.goos, #arch: flags.arch}
	matrix: {
		_goos: ["linux", "darwin"]
		_arch: ["amd64", "arm64"]
		for _g in _goos {
			// local, by goos
			(_g): _maker & {#variant: _g, #goos: _g, #arch: flags.arch}
			// arch x goos
			for _a in _arch {
				"\(_g)-\(_a)": _maker & {#variant: "\(_g)-\(_a)", #goos: _g, #arch: _a}
			}
		}
	}
}
