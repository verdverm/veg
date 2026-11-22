package runtime

import (
	"time"

	"github.com/verdverm/veg/pkg/flow/flow"
	"github.com/verdverm/veg/pkg/veg"
)

type FlowEnricher func(*Runtime, *flow.Flow) error

func (R *Runtime) EnrichFlows(flows []string, enrich FlowEnricher) error {
	// fmt.Println("Flow.Runtime.Enrich", flows)
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.Add("enrich/flow", end.Sub(start))
	}()

	if R.Flags.Verbosity > 1 {
		//fmt.Println("Runtime.Flow: ", flows)
		//for _, node := range R.Nodes {
		//  node.Print()
		//}
	}

	// Find only the datamodel nodes
	// TODO, dedup any references
	fs := []*flow.Flow{}
	for _, node := range R.Nodes {
		// check for Chat root
		if node.Veg.Flow.Root {
			// fmt.Println("flow.root:", node.Veg.Path)
			if !keepFilter(node, flows) {
				continue
			}
			upgrade := func(n *veg.Node[flow.Flow]) *flow.Flow {
				v := flow.NewFlow(n)
				return v
			}
			u := veg.Upgrade[any, flow.Flow](node, upgrade, nil)
			// we'd like this line in upgrade, but...
			// how do we make T a Node[T] type (or ensure that it has a veg)
			// u.T.Veg = u.Veg
			f := u.T
			f.Node = u
			fs = append(fs, f)
		}
	}

	R.Workflows = fs

	for _, c := range R.Workflows {
		err := enrich(R, c)
		if err != nil {
			return err
		}
	}


	return nil
}
