package st

import (
	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
	"github.com/verdverm/veg/pkg/structural"
)

type Upsert struct{}

func NewUpsert(val cue.Value) (vegcontext.Runner, error) {
	return &Upsert{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Upsert) Run(ctx *vegcontext.Context) (interface{}, error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()
	v := ctx.Value

	x := v.LookupPath(cue.ParsePath("val"))
	u := v.LookupPath(cue.ParsePath("upsert"))

	r, err := structural.UpsertValue(u, x, nil)
	if err != nil {
		return nil, err
	}

	r, err = v.FillPath(cue.ParsePath("out"), r), nil
	return r, err
}
