package cmds

import (
	"github.com/verdverm/veg/catalogs/gen/cli/schema"
)

TuiCommand: schema.Command & {
	Name:  "tui"
	Usage: "tui"
	Short: "a terminal interface to Hof and CUE"
	Long:  Short
	// Hidden: true
}
