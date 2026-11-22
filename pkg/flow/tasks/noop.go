package tasks

import (
	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type Noop struct{}

func NewNoop(val cue.Value) (vegcontext.Runner, error) {
	return &Noop{}, nil
}

func (T *Noop) Run(ctx *vegcontext.Context) (interface{}, error) {
	return nil, nil
}
