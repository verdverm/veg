package st

import (
	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
	"github.com/verdverm/veg/pkg/structural"
)

type Pick struct{}

func NewPick(val cue.Value) (vegcontext.Runner, error) {
	return &Pick{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Pick) Run(ctx *vegcontext.Context) (interface{}, error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()

	v := ctx.Value

	x := v.LookupPath(cue.ParsePath("val"))
	p := v.LookupPath(cue.ParsePath("pick"))

	r, err := structural.PickValue(p, x, nil)
	if err != nil {
		return nil, err
	}

	return v.FillPath(cue.ParsePath("out"), r), nil
}
