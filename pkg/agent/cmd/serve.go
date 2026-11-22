package cmd

import (
	"github.com/verdverm/veg/cmd/veg/flags"
	"github.com/verdverm/veg/pkg/yagu"
	"github.com/olekukonko/tablewriter"
)

func Serve(args []string, rflags flags.RootPflagpole, aflags flags.AgentPflagpole) error {
	R, _, matches, err := commonStart(args, rflags, aflags)
	if err != nil {
		return err
	}

	// gather rows
	var rows = make([][]string, 0, len(R.Agentics))
	// fill with data
	for _, e := range matches {
		_, kind, mname := extractMeta(e)

		path := e.Veg.Path
		// extra := genExtra(e)

		row := []string{mname, kind, path}
		rows = append(rows, row)
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Kind", "Path"},
		func(table *tablewriter.Table) ([][]string, error) {
			return rows, nil
		},
	)
}
