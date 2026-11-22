package tasker

import (
	"fmt"

	"cuelang.org/go/cue"
	cueflow "cuelang.org/go/tools/flow"

	flowctx "github.com/verdverm/veg/pkg/flow/context"
	"github.com/verdverm/veg/pkg/flow/task"
	"github.com/verdverm/veg/pkg/cuetils"
	"github.com/verdverm/veg/pkg/veg"
)

var debug = false

func NewTasker(ctx *flowctx.Context) cueflow.TaskFunc {
	// This function implements the Runner interface.
	// It parses Cue values, you will see all of them recursively

	return func(val cue.Value) (cueflow.Runner, error) {
		// fmt.Println("Tasker:", val.Path())

		// what's going on here?
		// this the root value? (so def a flow, so not a task)
		if len(val.Path().Selectors()) == 0 {
			return nil, nil
		}

		node, err := veg.ParseVeg[any](val)
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, nil
		}
		if node.Veg.Flow.Task == "" {
			return nil, nil
		}
		//if node.Veg.Flow.Task == "nest" {
		//  fmt.Println("New Tasker NEST", node.Veg.Path, node.Veg.Label)
		//}

		return makeTask(ctx, node)
	}
}

func makeTask(ctx *flowctx.Context, node *veg.Node[any]) (cueflow.Runner, error) {

	taskId := node.Veg.Flow.Task

	//taskName := node.Veg.Flow.Name
	//fmt.Println("makeTask:", taskId, taskName, node.Veg.Path, node.Veg.Flow.Root, node.Parent)

	// lookup context.RunnerFunc
	runnerFunc := ctx.Lookup(taskId)
	if runnerFunc == nil {
		return nil, fmt.Errorf("unknown task: %q at %q", taskId, node.Value.Path())
	}

	// Note, we apply this in the reverse order so that the Use order is like a stack
	// (i.e. the first is the most outer, which is typical for how these work for servers
	// apply plugin / middleware
	for i := len(ctx.Middlewares) - 1; i >= 0; i-- {
		ware := ctx.Middlewares[i]
		runnerFunc = ware.Apply(ctx, runnerFunc)
	}

	// some way to validate task against it's schema
	// (1) schemas self register
	// (2) here, we lookup schemas by taskId
	// (3) use custom Require (or other validator)

	// create veg task from val
	// these live under /flow/tasks
	// and are of type context.RunnerFunc
	T, err := runnerFunc(node.Value)
	if err != nil {
		return nil, err
	}

	// do per-task setup / common base / initial value / bookkeeping
	bt := task.NewBaseTask(node)
	ctx.Tasks.Store(bt.ID, bt)

	// wrap our RunnerFunc with cue/flow RunnerFunc
	return cueflow.RunnerFunc(func(t *cueflow.Task) error {
		// fmt.Println("makeTask.func()", t.Index(), t.Path())

		// why do we need a copy?
		// maybe for local Value / CurrTask
		c := flowctx.Copy(ctx)

		c.Value = t.Value()
		node, err := veg.ParseVeg[any](c.Value)
		if err != nil {
			return err
		}
		if node.Veg.Flow.Print.Level > 0 && node.Veg.Flow.Print.Before {
			pv := c.Value.LookupPath(cue.ParsePath(node.Veg.Flow.Print.Path))
			if node.Veg.Path == "" {
				fmt.Printf("%s", node.Veg.Flow.Print.Path)
			} else if node.Veg.Flow.Print.Path == "" {
				fmt.Printf("%s", node.Veg.Path)
			} else {
				fmt.Printf("%s.%s", node.Veg.Path, node.Veg.Flow.Print.Path)
			}
			fmt.Printf(": %#v\n", pv)
		}
		
		c.BaseTask = bt

		// fmt.Println("MAKETASK", taskId, c.FlowStack, c.Value.Path())
		// fmt.Printf("%# v\n", c.Value)

		bt.CueTask = t
		bt.Start = c.Value
		// TODO, we should remove this next line, and only set Final at the end
		bt.Final = c.Value

		// run the veg task
		bt.AddTimeEvent("run.beg")
		// (update)
		value, rerr := T.Run(c)
		bt.AddTimeEvent("run.end")

		if value != nil {
			// fmt.Println("FILL:", taskId, c.Value.Path(), t.Value(), value)
			bt.AddTimeEvent("fill.beg")

			//if node.Veg.Flow.Task == "nest" || node.Veg.Flow.Task == "api.Call" {
			//  fmt.Println("FILL:", taskId, c.Value.Path(), value)
			//}
			err = t.Fill(value)
			bt.Final = t.Value()
			bt.AddTimeEvent("fill.end")

			// fmt.Println("FILL:", taskId, c.Value.Path(), t.Value(), value)
			if err != nil {
				c.Error = err
				bt.Error = err
				return err
			}

			//if node.Veg.Flow.Print.Level > 0 && !node.Veg.Flow.Print.Before {
			//  // pv := bt.Final.LookupPath(cue.ParsePath(node.Veg.Flow.Print.Path))
			//  fmt.Printf("%s.%s: %# v\n", node.Veg.Path, node.Veg.Flow.Print.Path, value)
			//}

		}

		if rerr != nil {
			rerr = fmt.Errorf("in %q\n%v\n%+v", c.Value.Path(), cuetils.ExpandCueError(rerr), value)
			// fmt.Println("RunnerRunc Error:", err)
			c.Error = rerr
			bt.Error = rerr
			return rerr
		}

		return nil
	}), nil
}
