package info

import (
	"cuelang.org/go/cue"

	"github.com/verdverm/veg/cmd/veg/flags"
	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type BookkeepingConfig struct {
	Workdir string
}

type Bookkeeping struct {
	cfg  BookkeepingConfig
	val  cue.Value
	next vegcontext.Runner
}

func NewBookkeeping(cfg BookkeepingConfig, opts flags.RootPflagpole, popts flags.FlowPflagpole) *Bookkeeping {
	return &Bookkeeping{
		cfg: cfg,
	}
}

func (M *Bookkeeping) Run(ctx *vegcontext.Context) (results interface{}, err error) {
	// bt := ctx.BaseTask
	// fmt.Println("bt:", bt.ID, bt.UUID)
	result, err := M.next.Run(ctx)

	// write out file in background
	return result, err
}

func (M *Bookkeeping) Apply(ctx *vegcontext.Context, runner vegcontext.RunnerFunc) vegcontext.RunnerFunc {
	return func(val cue.Value) (vegcontext.Runner, error) {
		// id := fmt.Sprint(val.Path())
		// fmt.Println("book: found @", val.Path())
		next, err := runner(val)
		if err != nil {
			return nil, err
		}
		return &Bookkeeping{
			val:  val,
			next: next,
		}, nil
	}
}

func (M *Bookkeeping) write(filename string, val cue.Value) error {

	return nil
}
