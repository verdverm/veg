package cmd

import (
	// "fmt"

	"github.com/codemodus/kace"
	"github.com/olekukonko/tablewriter"

	"github.com/verdverm/veg/cmd/veg/flags"
	"github.com/verdverm/veg/pkg/datamodel"
	"github.com/verdverm/veg/pkg/runtime"
	"github.com/verdverm/veg/pkg/yagu"
)

func list(R *runtime.Runtime, dflags flags.DatamodelPflagpole) error {
	return yagu.PrintAsTable(
		[]string{"Name", "Type", "Version", "Status", "ID"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(R.Datamodels))
			// fill with data
			for _, dm := range R.Datamodels {
				id := dm.Veg.ID
				if id == "" {
					id = kace.Snake(dm.Veg.Metadata.Name) + " (auto)"
				}

				name := dm.Veg.Metadata.Name
				typ := datamodel.DatamodelType(dm)
				ver := dm.Veg.Datamodel.Version
				if ver == "" {
					ver = "-"
				}
				status := dm.Status()
				if status == "" {
					status = "-"
				}

				row := []string{name, typ, ver, status, id}
				rows = append(rows, row)
			}
			return rows, nil
		},
	)
}
