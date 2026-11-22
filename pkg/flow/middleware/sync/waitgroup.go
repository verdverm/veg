package sync

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/verdverm/veg/cmd/veg/flags"
	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type Wait struct {
	// "required"
	val  cue.Value
	next vegcontext.Runner

	// local
	name string
}

func NewWait(opts flags.RootPflagpole, popts flags.FlowPflagpole) *Wait {
	return &Wait{}
}

func (M *Wait) Run(ctx *vegcontext.Context) (results interface{}, err error) {
	// wg := ctx.<middleware>.[M.name]
	// wg.Add(1)
	// defer wg.Done()

	fmt.Println("wait: pre @", M.val.Path())
	// should this happen during discovery? (in Apply)
	result, err := M.next.Run(ctx)
	fmt.Println("wait: post @", M.val.Path())

	return result, err
}

func (M *Wait) Apply(ctx *vegcontext.Context, runner vegcontext.RunnerFunc) vegcontext.RunnerFunc {
	return func(val cue.Value) (vegcontext.Runner, error) {
		hasAttr := false
		attrs := val.Attributes(cue.ValueAttr)
		var a cue.Attribute
		for _, attr := range attrs {
			if attr.Name() == "wait" {
				a = attr
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

		fmt.Println("wait: found @", val.Path(), a)

		// what is in this attribute

		// setup wait by name here

		return &Wait{
			// required
			val:  val,
			next: next,
			// extra
			// name:
		}, nil
	}
}
