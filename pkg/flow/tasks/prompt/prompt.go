package prompt

import (
	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
	libprompt "github.com/verdverm/veg/pkg/prompt"
)

type Prompt struct{}

func NewPrompt(val cue.Value) (vegcontext.Runner, error) {
	return &Prompt{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Prompt) Run(ctx *vegcontext.Context) (any, error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()

	r, err := libprompt.RunPrompt(ctx.Value)
	if err != nil {
		return nil, err
	}

	return r, nil
}
