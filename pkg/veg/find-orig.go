package veg

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/codemodus/kace"
)

func FindVegsOrig(value cue.Value) (roots []*Node[any], err error) {
	// fmt.Println("FindVegs!")
	var stack *Node[any] // cue stack
	var nodes *Node[any] // veg nodes

	before := func(val cue.Value) bool {
		// get some info
		path := val.Path()
		sels := path.Selectors()
		last := cue.Selector{}
		label := ""
		if len(sels) > 0 {
			last = sels[len(sels)-1]
			label = last.String()
		}

		// return early and recurse for root value
		if label == "" {
			label = "<<root>>"
			// return true
		}

		// do not decend into $veg value itself
		// or any definition
		if label == "#veg" {
			return false
		}
		if label == "History" {
			return false
		}

		// update cue stack
		curr := New(label, val, nil, stack)
		stack = curr

		// did we find something of interest?
		found := false

		// look for #veg: _
		hv := val.LookupPath(cue.ParsePath("#veg"))
		if hv.Exists() {
			found = true
			err = hv.Decode(&(stack.Veg))
			if err != nil {
				fmt.Println(err)
				stack = stack.Parent
				return false
			}

			// fmt.Println("id:", stack.Veg.Label, stack.Veg.Path, stack.Veg.Metadata)
		}

		//
		// Do not modify stack.Veg before we Decode it in the condition just above
		//

		// need to re-add label, path, name? here, after decode
		stack.Veg.Label = label
		stack.Veg.Path = val.Path().String()

		ufound := upgradeAttrs(stack, label)
		if ufound {
			found = ufound
		}

		// filters to end recursion
		// check datamodel root because of nested history and roots snafu
		if stack.Veg.Datamodel.Root {
			// backtrack, walking parents
			for bt := nodes; bt != nil; bt = bt.Parent {
				// we found a nested root datamodel
				if bt.Veg.Datamodel.Root {
					// stop recursion
					fmt.Println("veg.DM: want to stop recursion here", bt.Veg.Path, stack.Veg.Path)
					// return false
				}
			}

			// fmt.Println("found datamodel:", stack.Veg.Path)
		}
		//if stack.Veg.Label == "History" {
		//  stack = stack.Parent
		//  return false
		//}

		// we should update the nodes
		if found {
			// update veg node
			node := New[any](label, val, nil, nodes)
			nodes = node
			nodes.Veg = stack.Veg

			// more enrichment
			if nodes.Veg.Metadata.Name == "" {
				nodes.Veg.Metadata.Name = nodes.Veg.Label
				nodes.Value = nodes.Value.FillPath(cue.ParsePath("#veg.metadata.name"), nodes.Veg.Metadata.Name)
			}
			if nodes.Veg.ID == "" {
				nodes.Veg.ID = nodes.Veg.Metadata.Name
				nodes.Value = nodes.Value.FillPath(cue.ParsePath("#veg.id"), kace.Kebab(nodes.Veg.ID))
			}

			if nodes.Parent == nil {
				// add to root if no parent
				roots = append(roots, nodes)
			} else {
				// add to parent's Children
				nodes.Parent.Children = append(nodes.Parent.Children, nodes)
			}
		}

		return true
	}

	after := func(val cue.Value) {
		// paths for matching trees
		np, sp := "", ""
		if nodes != nil {
			np = nodes.Value.Path().String()
		}
		if stack != nil {
			sp = stack.Value.Path().String()
		}

		if nodes != nil && np == sp {
			// unwind veg nodes
			nodes = nodes.Parent
		}
		// unwind node stack
		if stack != nil {
			stack = stack.Parent
		}
	}

	// this is a depth first walk
	WalkOrig(value, before, after)

	return roots, nil
}

// Walk is an alternative to cue.Value.Walk which handles more field types
// You can customize this with your own options
// returning false will stop recursion for that node
func WalkOrig(v cue.Value, before func(cue.Value) bool, after func(cue.Value), options ...cue.Option) {

	// call before and possibly stop recursion
	if before != nil && !before(v) {
		return
	}

	// possibly recurse
	switch v.IncompleteKind() {
	case cue.StructKind:
		if options == nil {
			options = defaultWalkOptions
		}
		s, _ := v.Fields(options...)

		for s.Next() {
			Walk(s.Value(), before, after, options...)
		}

	case cue.ListKind:
		l, _ := v.List()
		for l.Next() {
			Walk(l.Value(), before, after, options...)
		}

		// no default (basic lit types)

	}

	if after != nil {
		after(v)
	}

}
