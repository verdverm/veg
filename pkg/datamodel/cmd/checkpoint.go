package cmd

import (
	"fmt"
	"regexp"
	"time"

	"github.com/verdverm/veg/cmd/veg/flags"
	"github.com/verdverm/veg/pkg/datamodel"
	"github.com/verdverm/veg/pkg/runtime"
)

func checkpoint(R *runtime.Runtime, dflags flags.DatamodelPflagpole, cflags flags.Datamodel__CheckpointFlagpole) error {
	// check suffix
	suffix := cflags.Suffix
	matched, _ := regexp.MatchString("^[a-z_]+$", suffix)
	if !matched {
		return fmt.Errorf("suffix must contain only lowercase and underscores")
	}

	timestamp := time.Now().UTC().Format(datamodel.CheckpointTimeFmt)
	fmt.Printf("creating checkpoint: %s_%s %q\n", timestamp, suffix, cflags.Message)

	for _, dm := range R.Datamodels {
		err := dm.MakeSnapshot(timestamp, dflags, cflags)
		if err != nil {
			return err
		}
	}

	return nil
}
