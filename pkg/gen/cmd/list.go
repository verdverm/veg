package cmd

import (
	"fmt"

	"github.com/codemodus/kace"
	"github.com/olekukonko/tablewriter"

	"github.com/verdverm/veg/cmd/veg/flags"
	"github.com/verdverm/veg/pkg/yagu"
)

func List(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole) error {
	R, err := prepRuntime(args, rflags, gflags)
	if err != nil {
		return err
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Path", "ID", "Creator"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(R.Generators))
			// fill with data
			for _, gen := range R.Generators {
				id := gen.Veg.ID
				if id == "" {
					id = kace.Snake(gen.Veg.Metadata.Name) + " (auto)"
				}

				name := gen.Veg.Gen.Name
				if name == "" {
					name = "(anon)"
				}
				path := gen.Veg.Path

				row := []string{name, path, id, fmt.Sprint(gen.Veg.Gen.Creator)}
				rows = append(rows, row)
			}
			return rows, nil
		},
	)
}
