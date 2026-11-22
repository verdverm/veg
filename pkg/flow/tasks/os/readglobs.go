package os

import (
	g_os "os"

	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
	"github.com/verdverm/veg/pkg/yagu"
)

type ReadGlobs struct{}

func NewReadGlobs(val cue.Value) (vegcontext.Runner, error) {
	return &ReadGlobs{}, nil
}

func (T *ReadGlobs) Run(ctx *vegcontext.Context) (interface{}, error) {

	val := ctx.Value

	patterns, err := extractGlobConfig(ctx, val)
	if err != nil {
		return nil, err
	}

	filepaths, err := yagu.FilepathsFromGlobs(patterns)
	if err != nil {
		return nil, err
	}

	data := map[string]string{}

	for _, file := range filepaths {
		bs, err := g_os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		data[file] = string(bs)
	}

	return map[string]interface{}{"files": data}, nil
}
