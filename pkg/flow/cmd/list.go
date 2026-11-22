package cmd

import (
	"github.com/codemodus/kace"
	"github.com/olekukonko/tablewriter"

	"github.com/verdverm/veg/cmd/veg/flags"
	"github.com/verdverm/veg/pkg/yagu"
)

func List(args []string, rflags flags.RootPflagpole, cflags flags.FlowPflagpole) error {
	R, err := prepRuntime(args, rflags, cflags)
	if err != nil {
		return err
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Path", "ID"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(R.Workflows))
			// fill with data
			for _, wf := range R.Workflows {
				id := wf.Veg.ID
				if id == "" {
					id = kace.Snake(wf.Veg.Metadata.Name) + " (auto)"
				}

				name := wf.Veg.Flow.Name
				if name == "" {
					name = "(anon)"
				}
				path := wf.Veg.Path

				row := []string{name, path, id}
				rows = append(rows, row)
			}
			return rows, nil
		},
	)
}
