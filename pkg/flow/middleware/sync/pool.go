package sync

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/gammazero/workerpool"

	"github.com/verdverm/veg/cmd/veg/flags"
	flowctx "github.com/verdverm/veg/pkg/flow/context"
	"github.com/verdverm/veg/pkg/veg"
)

type Pool struct {
	val  cue.Value
	next flowctx.Runner
}

func NewPool(opts flags.RootPflagpole, popts flags.FlowPflagpole) *Pool {
	// fmt.Println("Pool: new")
	return &Pool{}
}

func (M *Pool) Run(ctx *flowctx.Context) (results interface{}, err error) {
	val := ctx.Value

	node, err := veg.ParseVeg[any](val)
	if err != nil {
		return nil, err
	}
	if node == nil  {
		panic("we should have found a node to even get here")
	}

	vegp := node.Veg.Flow.Pool

	// Make (pool) is setup before hand? (as middleware?)
	if vegp.Take {

		pool, ok := ctx.Pools.Load(vegp.Name)
		if !ok {
			return nil, fmt.Errorf("unknown exec pool %q @ %s\n", vegp.Name, val.Path())
		}

		P, ok := pool.(*workerpool.WorkerPool)
		P.SubmitWait(func() {
			results, err = M.next.Run(ctx)
		})
	} else {
		// fmt.Println("Pool: skip @", M.val.Path())
		results, err = M.next.Run(ctx)
	}

	return results, err
}

func (M *Pool) Apply(ctx *flowctx.Context, runner flowctx.RunnerFunc) flowctx.RunnerFunc {
	// fmt.Println("pool.Apply call")
	return func(val cue.Value) (flowctx.Runner, error) {

		// fmt.Println("pool.Apply func")

		// parse out the local #veg node data
		node, err := veg.ParseVeg[any](val)
		if err != nil {
			return nil, err
		}
		if node == nil  {
			panic("we should have found a node to even get here")
		}

		// convenience
		vegp := node.Veg.Flow.Pool

		// hmmm, not sure if this actually runs it?
		// probably does, which is why it is paired with a no-op
		// so that a task is found, maybe this is where we are doubling up?
		next, err := runner(val)
		if err != nil {
			return nil, err
		}

		// if not a pool making task, return
		if vegp.Name == "" {
			return next, nil
		}

		// fmt.Printf("Pool: found @ %s %#v\n", val.Path(), vegp)

		// setup pool by name here

		if vegp.Make {
			pool := workerpool.New(vegp.Number)
			// fmt.Println("Pool: store @", val.Path(), pool)
			ctx.Pools.Store(vegp.Name, pool)
		}

		return &Pool{
			// required
			val:  val,
			next: next,
			// extra
		}, nil
	}
}
