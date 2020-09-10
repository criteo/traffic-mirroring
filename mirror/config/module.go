package config

import (
	"encoding/json"
	"fmt"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

var moduleIndex = map[string]int{}

type Module struct {
	Type   string
	Name   string
	Config json.RawMessage
}

func CreateModule(b []byte) (mirror.Module, error) {
	mc := Module{}
	err := json.Unmarshal(b, &mc)
	if err != nil {
		return nil, fmt.Errorf("error reading module config: %w", err)
	}

	name := mc.Name
	if name == "" {
		name = fmt.Sprintf("%s.%d", mc.Type, moduleIndex[mc.Type])
	}
	moduleIndex[mc.Type]++

	ctx := mirror.ModuleContext{
		Name: name,
	}

	return registry.Create(mc.Type, ctx, mc.Config)
}

func CreateModules(b []byte) ([]mirror.Module, error) {
	res := []mirror.Module{}

	c := []json.RawMessage{}
	err := json.Unmarshal(b, &c)
	if err != nil {
		return nil, fmt.Errorf("error reading module config: %w", err)
	}

	for _, mc := range c {
		mod, err := CreateModule(mc)
		if err != nil {
			return nil, err
		}

		res = append(res, mod)
	}

	return res, nil
}
