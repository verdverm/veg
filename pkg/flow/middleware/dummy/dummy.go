package dummy

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/verdverm/veg/cmd/veg/flags"
	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type Dummy struct {
	val  cue.Value
	next vegcontext.Runner
}

func NewDummy(opts flags.RootPflagpole, popts flags.FlowPflagpole) *Dummy {
	return &Dummy{}
}

func (M *Dummy) Run(ctx *vegcontext.Context) (results interface{}, err error) {
	fmt.Println("dummy: pre @", M.val.Path())
	// should this happen during discovery? (in Apply)
	result, err := M.next.Run(ctx)
	fmt.Println("dummy: post @", M.val.Path())

	return result, err
}

func (M *Dummy) Apply(ctx *vegcontext.Context, runner vegcontext.RunnerFunc) vegcontext.RunnerFunc {
	return func(val cue.Value) (vegcontext.Runner, error) {
		hasAttr := false
		attrs := val.Attributes(cue.ValueAttr)
		for _, attr := range attrs {
			if attr.Name() == "dummy" {
				hasAttr = true
				break
			}
		}

		next, err := runner(val)
		if err != nil {
			return nil, err
		}

		if !hasAttr {
			return next, nil
		}

		fmt.Println("dummy: found @", val.Path())

		return &Dummy{
			val:  val,
			next: next,
		}, nil
	}
}
