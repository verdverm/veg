package env

import "github.com/verdverm/veg/pkg/veg"

type Env struct {
	*veg.Node[any]

	Name        string
	HumanName   string
	MachineName string

	Kind string

	Description        string
	HumanDescription   string
	MachineDescription string

	// stuff to match docker / compose / k8s
}
