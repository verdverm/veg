@experiment(aliasv2)
package config

import (
	"strings"

	"github.com/verdverm/veg/schemas/env"
)

let root = self

bins: multi: [string]: _

dist: {
	[!~"(images|sbom)"]~(k,_): {@env()
		#veg: {id: "dist-\(k)", metadata: {name: string | *id}}
		name: string | *#veg.metadata.name
	}

	meta: env.#ExportDir & {
		path: "dist/meta"
		sources: [
			// root.src.changelog,
			// bins.checksum,

			for sbom in dist.sbom {sbom},
		]
		wipe: true
	}

	cuemod: env.#ExportDir & {
		name: "cue-module"
		path: "dist/cuemod"
		sources: [root.src.cuemod]
		wipe: true
	}

	github: env.#ExportDir & {
		name: "github-files"
		path: "dist/github"
		sources: [
			for key, val in veg.cli.matrix if key != "name" {val},

		]
		trimPrefix: "./bins"
		wipe:       true
	}

	vscode: env.#ExportDir & {
		name: "vscode-files"
		path: "dist/vscode"
		sources: [
			extn.vscode.vsix,
		]
		wipe: true
	}

	images: {
		[string]~(k,_): env.#PublishImage & {
			@env()
			#veg: id: string | *"dist-veg-\(k)"
			#veg: metadata: name: #veg.id
			reg:  root.flags.registry
			name: string | *"veg-\(k)"
		}
		min: {image: root.ctr.min}
		dev: {image: root.ctr.dev}
		// the dev image with adk/dagger added
		veg: {
			image: env.#Container & {
				from: root.ctr.dev
				steps: [
					env.Dir & {path: "/adk", source: src.adk.fork},
					env.Dir & {path: "/dagger", source: src.dagger.fork},
				]
			}
		}
		ops: {image: root.ctr["ops-all"], name: "veg-ops"}
		for f, F in root.fmtr {
			let _f = "fmt-\(f)"
			(_f): {
				@env()
				#veg: id: "dist-\(_f)"
				#veg: metadata: name: #veg.id
				name: "\(_f)"
				image: F.img
			}
		}
	}

	sbom: {
		for k in ["cuemod", "github", "vscode"] {
			(k): env.#CuefigSBOM & {
				@env()
				#veg: id: "sbom-\(k)"
				#veg: metadata: name: #veg.id
				path:   "\(k).cue"
				format: "cue"
				data:   dist[k]
			}
		}
		for i, img in dist.images {
			(i): env.#CuefigSBOM & {
				@env()
				_id: strings.TrimPrefix(img.#veg.id, "dist-")
				#veg: id: "sbom-\(_id)"
				#veg: metadata: name: #veg.id
				path:   "\(_id).cue"
				format: "cue"
				data:   img
			}
		}
	}

}
