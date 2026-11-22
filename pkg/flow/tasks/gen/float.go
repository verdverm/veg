package gen

import (
	"math/rand"

	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type Float struct{}

func NewFloat(val cue.Value) (vegcontext.Runner, error) {
	return &Float{}, nil
}

func (T *Float) Run(ctx *vegcontext.Context) (interface{}, error) {
	f := rand.Float64()
	return f, nil
}

type Norm struct{}

func NewNorm(val cue.Value) (vegcontext.Runner, error) {
	return &Norm{}, nil
}

func (T *Norm) Run(ctx *vegcontext.Context) (interface{}, error) {
	f := rand.NormFloat64()
	return f, nil
}
