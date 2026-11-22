package veg

import (
	"github.com/verdverm/veg/catalogs/gen/cli/gen"

	"github.com/verdverm/veg/design"
)

Cli: gen.Generator & {
	@gen(cli,veg)
	Outdir:     "./"
	Cli:        design.CLI
	WatchGlobs: ["./design/**/*"]
	WatchExec: {
		@task(os.Exec)
		cmd: ["go", "install", "\(Outdir)/cmd/veg"]
		env: {
			CGO_ENABLE: "0"
		}
		exitcode: _
	}
}
