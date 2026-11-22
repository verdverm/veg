package gen

import (
	"cuelang.org/go/cue"
	"github.com/google/uuid"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type UUID struct{}

func NewUUID(val cue.Value) (vegcontext.Runner, error) {
	return &UUID{}, nil
}

func (T *UUID) Run(ctx *vegcontext.Context) (interface{}, error) {
	u := uuid.New()
	return u, nil
}
