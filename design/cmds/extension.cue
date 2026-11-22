package cmds

import (
	"github.com/verdverm/veg/catalogs/gen/cli/schema"
)

ExtensionCommand: schema.Command & {
	Name:  "extension"
	Usage: "extension [args]"
	Short: "run the extension server"
	Long:  "run the extension server"

	Imports: [
		{Path: "github.com/verdverm/veg/cmd/veg/flags"},
		{Path: "github.com/verdverm/veg/pkg/extension"},
		{Path: "github.com/verdverm/veg/pkg/runtime"},
	]

	PersistentPrerun:     true
	PersistentPrerunBody: "err = runtime.EnsureInfra()"

	// this runs the commands or naked `veg env` command
	Body: "err = extension.Run(args, flags.RootPflags)"
}
