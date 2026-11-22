package common

import (
	"github.com/verdverm/veg/pkg/agent/config"
	"github.com/verdverm/veg/pkg/agent/services/environ"
	"github.com/verdverm/veg/pkg/cuetils"
)

func ReloadConfig(ar Runtime) error {
	err := ar.ReadEnvConfig()
	if err != nil {
		return cuetils.ExpandCueError(err)
	}
	return nil
}

func ListEnvirons() ([]environ.Environ, error) {
	return environ.Client().ListEnvirons()
}

func GetModels(ar Runtime) (map[string]config.Model, error) {
	return ar.GetAgenticConfig().Models, nil
}

func GetAgents(ar Runtime) (map[string]config.Agent, error) {
	return ar.GetAgenticConfig().Agents, nil
}
