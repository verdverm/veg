package os

import (
	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
	"github.com/verdverm/veg/pkg/yagu"
)

type Glob struct{}

func NewGlob(val cue.Value) (vegcontext.Runner, error) {
	return &Glob{}, nil
}

func (T *Glob) Run(ctx *vegcontext.Context) (interface{}, error) {

	val := ctx.Value

	patterns, err := extractGlobConfig(ctx, val)
	if err != nil {
		return nil, err
	}

	filepaths, err := yagu.FilepathsFromGlobs(patterns)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"filepaths": filepaths}, nil
}

func extractGlobConfig(ctx *vegcontext.Context, val cue.Value) (patterns []string, err error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()

	ps := val.LookupPath(cue.ParsePath("globs"))
	if ps.Err() != nil {
		return nil, ps.Err()
	}

	iter, err := ps.List()
	if err != nil {
		return nil, err
	}

	for iter.Next() {
		gv := iter.Value()
		if gv.Err() != nil {
			return nil, gv.Err()
		}
		gs, err := gv.String()
		if err != nil {
			return nil, err
		}

		patterns = append(patterns, gs)
	}

	return patterns, nil
}
