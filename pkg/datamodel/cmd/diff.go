package cmd

import (
	"os"

	"github.com/verdverm/veg/cmd/veg/flags"
	"github.com/verdverm/veg/pkg/runtime"
)

func diff(R *runtime.Runtime, dflags flags.DatamodelPflagpole) error {

	for _, dm := range R.Datamodels {
		if !dm.HasDiff() {
			continue
		}
		if err := dm.PrintDiff(os.Stdout, dflags); err != nil {
			return err
		}
	}

	return nil
}
