package st

import (
	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
	"github.com/verdverm/veg/pkg/structural"
)

type Insert struct{}

func NewInsert(val cue.Value) (vegcontext.Runner, error) {
	return &Insert{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Insert) Run(ctx *vegcontext.Context) (interface{}, error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()

	v := ctx.Value

	x := v.LookupPath(cue.ParsePath("val"))
	ins := v.LookupPath(cue.ParsePath("insert"))

	r, err := structural.InsertValue(ins, x, nil)
	if err != nil {
		return nil, err
	}

	return v.FillPath(cue.ParsePath("out"), r), nil
}
