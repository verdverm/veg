package cmd

import (
	"github.com/verdverm/veg/cmd/veg/flags"
	"github.com/verdverm/veg/pkg/runtime"
)

func findMaxLabelLen(R *runtime.Runtime, dflags flags.DatamodelPflagpole) int {
	max := 0
	for _, dm := range R.Datamodels {
		m := dm.FindMaxLabelLen(dflags)
		if m > max {
			max = m
		}
	}
	return max
}