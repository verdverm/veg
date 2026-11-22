package gen

import (
	"cuelang.org/go/cue"
	"github.com/lucsky/cuid"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type CUID struct{}

func NewCUID(val cue.Value) (vegcontext.Runner, error) {
	return &CUID{}, nil
}

func (T *CUID) Run(ctx *vegcontext.Context) (interface{}, error) {
	u := cuid.New()
	return u, nil
}

type Slug struct{}

func NewSlug(val cue.Value) (vegcontext.Runner, error) {
	return &Slug{}, nil
}

func (T *Slug) Run(ctx *vegcontext.Context) (interface{}, error) {
	u := cuid.Slug()
	return u, nil
}
