package info

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/verdverm/veg/cmd/veg/flags"
	vegcontext "github.com/verdverm/veg/pkg/flow/context"
	"github.com/verdverm/veg/pkg/veg"
)

type Print struct {
	val  cue.Value
	next vegcontext.Runner
}

func NewPrint(opts flags.RootPflagpole, popts flags.FlowPflagpole) *Print {
	return &Print{}
}

func (M *Print) Run(ctx *vegcontext.Context) (results interface{}, err error) {
	result, err := M.next.Run(ctx)

	fmt.Println("HELLO")

	attrs := M.val.Attributes(cue.ValueAttr)
	for _, attr := range attrs {
		if attr.Name() == "print" {
			for i := 0; i < attr.NumArgs(); i++ {
				p, err := attr.String(i)
				if err != nil {
					return result, err
				}
				if p == "" {
					fmt.Printf("%s: %v\n", M.val.Path(), result)
				} else {
					r, ok := result.(cue.Value)
					var v cue.Value
					if ok {
						v = r.LookupPath(cue.ParsePath(p))
					} else {
						v = M.val.LookupPath(cue.ParsePath(p))
					}
					fmt.Printf("%s: %v\n", v.Path(), v)
				}
			}
			break
		}
	}

	return result, err
}

func (M *Print) Apply(ctx *vegcontext.Context, runner vegcontext.RunnerFunc) vegcontext.RunnerFunc {
	return func(val cue.Value) (vegcontext.Runner, error) {
		// parse out the local #veg node data
		node, err := veg.ParseVeg[any](val)
		if err != nil {
			return nil, err
		}
		if node == nil  {
			panic("we should have found a node to even get here")
		}

		p := node.Veg.Flow.Print

		// convenience
		fmt.Println("found print", val.Path(), node.Veg.Path, p)

		next, err := runner(val)
		if err != nil {
			return nil, err
		}

		if p.Level > 0 {
			return next, nil
		}

		fmt.Println("found print")

		return &Print{
			val:  val,
			next: next,
		}, nil
	}
}
