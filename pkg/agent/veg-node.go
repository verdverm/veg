package agent

import "github.com/verdverm/veg/pkg/veg"

type Agentic struct {
	*veg.Node[any]

	Name        string
	HumanName   string
	MachineName string

	Description        string
	HumanDescription   string
	MachineDescription string

	// stuff to match docker / compose / k8s
}
