package datamodel

import (
	"github.com/verdverm/veg/cmd/veg/flags"
)

func (dm *Datamodel) FindMaxLabelLen(dflags flags.DatamodelPflagpole) int {
	max := len(dm.Veg.Label)
	m := dm.T.findMaxLabelLenR("", "  ", dflags)
	if m > max {
		max = m
	}
	return max
}

func (V *Value) findMaxLabelLenR(indent, spaces string, dflags flags.DatamodelPflagpole) int {
	max := V.findMaxLabelLen(indent, spaces, dflags)
	for _, c := range V.Children {
		m := c.T.findMaxLabelLenR(indent + spaces, spaces, dflags)
		if m > max {
			max = m
		}
	}
	return max
}

func (V *Value) findMaxLabelLen(indent, spaces string, dflags flags.DatamodelPflagpole) int {
	return len(V.Veg.Label) + len(indent)
}
