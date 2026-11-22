package veg

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"github.com/codemodus/kace"
)

// this is where we upgrade our @sugar() to equivalent #veg: { ... } in the Golang Veg type
func upgradeAttrs[T any](node *Node[T], label string) bool {
	val := node.Value
	attrs := val.Attributes(cue.ValueAttr)

	// did we find any attribute
	found := false
	for _, A := range attrs {
		an, ac := A.Name(), A.Contents()
		ac = strings.TrimSpace(ac)
		lfound := true
		switch an {
		case "veg":
			switch ac {
			case "datamodel":
				node.Veg.Datamodel.Root = true
			}
		case "id":
			node.Veg.ID = ac
		case "memo":
			node.Veg.Memo = ac
		case "z":
			i, err := strconv.Atoi(ac)
			if err == nil {
				node.Veg.Z = i
			} else {
				fmt.Fprintf(os.Stderr, "WARN: @z(%s) parsing error in %s: %s", ac, val.Path(), err)
			}

		case "name":
			node.Veg.Metadata.Name = ac

		case "datamodel":
			node.Veg.Datamodel.Root = true
		case "history":
			node.Veg.Datamodel.History = true
		case "ordered":
			node.Veg.Datamodel.Ordered = true
		case "cue":
			node.Veg.Datamodel.Cue = true

		// doesn't handle empty case, do we support that
		// we probably should
		case "gen":
			node.Veg.Gen.Root = true
			node.Veg.Gen.Name = ac
			c := val.LookupPath(cue.ParsePath("Create"))
			if c.Exists() {
				node.Veg.Gen.Creator = true
			}

		// this doesnt handle empty @flow()
		case "flow":
			node.Veg.Flow.Root = true
			node.Veg.Flow.Name = ac
			// extra, a bit hacky
			node.Veg.Metadata.Name = ac

		// this doesn't handle task names
		// maybe we split into parts
		case "task":
			node.Veg.Flow.Task = ac
			node.Veg.Flow.Name = label

		case "pool":

			parts := strings.Split(ac, ",")

			if len(parts) > 1 {
				node.Veg.Flow.Pool.Name = parts[0]
				node.Veg.Flow.Pool.Make = true
				c, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println("warning: unable to parse %q to int", parts[1])

				} else {
					node.Veg.Flow.Pool.Number = c
				}

			} else {
				node.Veg.Flow.Pool.Name = ac
				node.Veg.Flow.Pool.Take = true
			}

		case "print":
			// TODO, better parsing of AC to get parts
			node.Veg.Flow.Print.Level = 1
			node.Veg.Flow.Print.Path = ac

		case "env":
			node.Veg.Env.Root = true
			// name preference
			// 0. @env(hide|hidden) will override any prior settings
			if ac == "hide" || ac == "hidden" || node.Veg.Env.Name == "hidden" {
				// this should hopefully overwrite, and give us order independence, hide always winning
				node.Veg.Env.Name = "hidden"
				break
			}

			// 1. name set manually
			// 2. @env(contents)
			// 3. @id(contents)
			// 4. label of the struct it is within
			n := val.LookupPath(cue.ParsePath("name"))
			if n.Exists() {
				if s, err := n.String(); err == nil {
					if s != "" {
						node.Veg.Env.Name = s
					}
				}
			}

			// wasn't manually set, lets do fallback
			if node.Veg.Env.Name == "" {
				if ac != "" {
					node.Veg.Env.Name = ac
				} else if node.Veg.ID != "" {
					node.Veg.Env.Name = node.Veg.ID
				} else {
					node.Veg.Env.Name = label
				}
			}

			// @env(...) will also write @id() @name() if not set already
			// trying this out, may use elsewhere, would be good to have a pattern
			if ac != "" {
				if node.Veg.ID == "" {
					node.Veg.ID = ac
				}
				if node.Veg.Metadata.Name == "" {
					node.Veg.Metadata.Name = ac
				}
			}
			// if not already, set metadata name based on env name
			if node.Veg.Metadata.Name == "" {
				node.Veg.Metadata.Name = node.Veg.Env.Name
			}

		case "agentic":
			node.Veg.Agentic.Root = true
			node.Veg.Agentic.Name = label
			node.Veg.Agentic.Kind = ac
			// node.Veg.Agentic.Extra = ac

		default:
			lfound = false
		}
		// write to outer found
		if lfound {
			found = true
			node.Veg.AtMade = true
		}
	}

	return found
}

// parse out a #veg for a single value
func ParseVeg[T any](val cue.Value) (*Node[T], error) {

	// get some info
	path := val.Path()
	sels := path.Selectors()
	last := cue.Selector{}
	label := ""
	if len(sels) > 0 {
		last = sels[len(sels)-1]
		label = last.String()
	}

	// do not decend into #veg value itself
	// or any definition
	if label == "#veg" {
		return nil, fmt.Errorf("you are trying to parse the #veg value, rather than the containing node at %v", path)
	}

	// return early and recurse for root value?
	if label == "" {
		label = "<<root>>"
		// return nil, nil
	}

	// create new node
	node := New[T](label, val, nil, nil)
	found := false

	// look for #veg: _
	hv := val.LookupPath(cue.ParsePath("#veg"))
	if hv.Exists() {
		err := hv.Decode(&(node.Veg))
		if err != nil {
			return node, err
		}
		found = true
	}

	// need to re-add label, path, name? here, after decode
	node.Veg.Label = label
	node.Veg.Path = path.String()

	//
	// look for attributes
	ufound := upgradeAttrs(node, label)
	if ufound {
		found = ufound
	}

	// filters to end recursion
	// check datamodel root because of nested history and roots snafu
	// This was here from before we made this singular function, in the context of nodes and the stack
	//if node.Veg.Datamodel.Root {
	//  // backtrack, walking parents
	//  for bt := nodes; bt != nil; bt = bt.Parent {
	//    // we found a nested root datamodel
	//    if bt.Veg.Datamodel.Root {
	//      // stop recursion
	//      fmt.Println("veg.DM: want to stop recursion here", bt.Veg.Path, node.Veg.Path)
	//      // return false
	//    }
	//  }
	//  // fmt.Println("found datamodel:", stack.Veg.Path)
	//}

	//// skip any injected History labels?
	//if node.Veg.Label == "History" {
	//  return nil, nil
	//}

	// HACK for nested task, so that it does not get picked up twice
	// (1) inside of a flow  (2) as the root of a flow
	// when we get to rethinking veg/flow with really using the #veg
	// we shouldn't even have something like pool or nest tasks at all
	// (also, CUE might support nested tasks more naturally)
	//if node.Veg.Flow.Pool.Take && node.Veg.Flow.Task == "nest" {
	//  return nil, nil
	//}

	if !found {
		return nil, nil
	}

	// These are EXPENSIVE, only do if we care about and will definitely be returning the node
	// more enrichment
	if node.Veg.Metadata.Name == "" {
		node.Veg.Metadata.Name = node.Veg.Label
		node.Value = node.Value.FillPath(cue.ParsePath("#veg.metadata.name"), node.Veg.Metadata.Name)
	}
	if node.Veg.ID == "" {
		node.Veg.ID = node.Veg.Metadata.Name
		node.Value = node.Value.FillPath(cue.ParsePath("#veg.metadata.id"), kace.Kebab(node.Veg.ID))
	}

	return node, nil
}

func FindVegs(value cue.Value) (roots []*Node[any], err error) {
	// fmt.Println("FindVegs!")
	var stack *Node[any] // cue stack

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

		if label == "#veg" {
			return false
		}

		node, err := ParseVeg[any](val)
		if err != nil {
			fmt.Println("while parsing for #veg:", err)
			return false
		}

		// did we find something?
		if node != nil {

			// skip any injected History labels?
			if node.Veg.Label == "History" {
				return false
			}

			// is this the root of that interesting thing?
			// otherwise, push onto the stack and update parent/child pointers
			if stack == nil {
				stack = node
				roots = append(roots, node)
			} else {
				// two-way relation setting
				node.Parent = stack
				stack.Children = append(stack.Children, node)

				// push stack
				stack = node
			}

			// hmm, should this go before or after the stack update?
			// it doesn't work if we do, so... trying to reconcile this with history
			// this is anything that can have a nested flow, ideally this could be inferred from the node / enrichment
			if node.Veg.Flow.Task == "nest" ||
				(node.Veg.Flow.Name != "" && node.Veg.Flow.Task == "") {
				// fmt.Println("ending veg recursion in nest", node.Veg.Path, node.Veg.Label)
				return false
			}

		}

		return true
	}

	after := func(val cue.Value) {
		// unwind node stack
		if stack != nil {
			stack = stack.Parent
		}
	}

	// this is a depth first walk
	Walk(value, before, after)

	return roots, nil
}

func (n *Node[T]) indent() string {
	if n == nil {
		return ""
	}
	d := 1
	for n.Parent != nil {
		d++
		n = n.Parent
	}
	return strings.Repeat("  ", d)
}

var defaultWalkOptions = []cue.Option{
	cue.Attributes(true),
	cue.Concrete(false),
	cue.Definitions(true),
	cue.Hidden(false),
	cue.Optional(true),
	cue.Docs(true),
}

// Walk is an alternative to cue.Value.Walk which handles more field types
// You can customize this with your own options
// returning false will stop recursion for that node
func Walk(v cue.Value, before func(cue.Value) bool, after func(cue.Value), options ...cue.Option) {

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
