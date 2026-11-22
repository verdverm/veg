package cmds

import (
	"github.com/verdverm/veg/catalogs/gen/cli/schema"
)

ModCmdImports: [
	{Path: "github.com/verdverm/veg/pkg/mod", ...},
	{Path: "github.com/verdverm/veg/cmd/veg/flags", ...},
]

ModCommand: schema.Command & {
	// TBD:   "β"
	Name:  "mod"
	Usage: "mod"
	Aliases: ["m"]
	Short: "CUE module dependency management"
	Long:  "CUE module dependency management, imported from upstream CUE project"

	//Topics: #ModTopics
	//Examples: #ModExamples
}
