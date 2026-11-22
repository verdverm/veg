package gen

import (
	"time"

	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type Now struct{}

func NewNow(val cue.Value) (vegcontext.Runner, error) {
	return &Now{}, nil
}

func (T *Now) Run(ctx *vegcontext.Context) (interface{}, error) {
	t := time.Now()
	f := t.Format(time.RFC3339)
	return f, nil
}
