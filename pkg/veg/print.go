package veg

import (
	"fmt"
)

func (hn *Node[T]) Print() {
	if hn.Veg.Datamodel.Root {
		hn.PrintDatamodel("  ", "")
	}

	if hn.Veg.Gen.Root {
		hn.PrintGen("  ", "")
	}

	if hn.Veg.Flow.Root {
		hn.PrintFlow("  ", "")
	}
}

func (hn *Node[T]) PrintDatamodel(spacer, indent string) {
	s := fmt.Sprintf("%s%s:", indent, hn.Veg.Path)
	dm := hn.Veg.Datamodel
	found := false
	if dm.Root {
		s += " @datamodel()"
		found = true
	}
	if dm.Node {
		s += " @node()"
		found = true
	}
	if dm.History {
		s += " @history()"
		found = true
	}
	if dm.Ordered {
		s += " @ordered()"
		found = true
	}
	if dm.Cue {
		s += " @cue()"
		found = true
	}
	
	if !found {
		return
	}
	// print line
	fmt.Println(s)
	// pretty.Println(hn.Veg)

	// recurse
	for _, c := range hn.Children {
		c.PrintDatamodel(spacer, indent+spacer)
	}
}

func (hn *Node[T]) PrintGen(spacer, indent string) {
	fmt.Printf("%s%s: @gen(%s)\n", indent, hn.Veg.Path, hn.Veg.Gen.Name)
	for _, c := range hn.Children {
		c.PrintDatamodel(spacer, indent+spacer)
	}
}

func (hn *Node[T]) PrintFlow(spacer, indent string) {
	f := hn.Veg.Flow
	s := fmt.Sprintf("%s%s:", indent, hn.Veg.Path)

	if f.Root {
		s += fmt.Sprintf(" @flow(%s)", f.Name)
	}
	if f.Task != "" {
		s += fmt.Sprintf(" @task(%s)", f.Task)
	}
	// print line
	fmt.Println(s)
	// pretty.Println(hn.Veg)

	// recurse
	for _, c := range hn.Children {
		c.PrintFlow(spacer, indent+spacer)
	}
}
