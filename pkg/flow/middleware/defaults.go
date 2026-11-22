package middleware

import (
	"github.com/verdverm/veg/cmd/veg/flags"
	vegcontext "github.com/verdverm/veg/pkg/flow/context"

	"github.com/verdverm/veg/pkg/flow/middleware/info"
	"github.com/verdverm/veg/pkg/flow/middleware/sync"
)

func UseDefaults(ctx *vegcontext.Context, opts flags.RootPflagpole, popts flags.FlowPflagpole) {
	// ctx.Use(dummy.NewDummy(opts, popts))
	// ctx.Use(info.NewPrint(opts, popts))
	ctx.Use(info.NewProgress(opts, popts))
	//ctx.Use(info.NewBookkeeping(info.BookkeepingConfig{
	//Workdir: ".veg/flow",
	//}, opts, popts))
	ctx.Use(sync.NewPool(opts, popts))
}
