package gen

import (
	"math/rand"
	"time"

	"cuelang.org/go/cue"

	vegcontext "github.com/verdverm/veg/pkg/flow/context"
)

type Seed struct{}

func NewSeed(val cue.Value) (vegcontext.Runner, error) {
	return &Seed{}, nil
}

func (T *Seed) Run(ctx *vegcontext.Context) (interface{}, error) {

	val := ctx.Value

	var s int64
	var err error

	ferr := func() error {
		ctx.CUELock.Lock()
		defer func() {
			ctx.CUELock.Unlock()
		}()

		// lookup key
		sv := val.LookupPath(cue.ParsePath("seed"))
		if sv.Exists() {
			if sv.Err() != nil {
				return sv.Err()
			}
			s, err = sv.Int64()
			if err != nil {
				return err
			}
		} else {
			s = time.Now().UnixNano()
		}

		return nil
	}()
	if ferr != nil {
		return nil, ferr
	}

	rand.Seed(s)

	return nil, nil
}
