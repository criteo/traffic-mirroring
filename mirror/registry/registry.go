package registry

import (
	"fmt"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
)

type FactoryFunc func(cfg []byte) (mirror.Module, error)

type Registry struct {
	modules map[string]FactoryFunc
}

func New() *Registry {
	return &Registry{
		modules: map[string]FactoryFunc{},
	}
}

func (r *Registry) Register(name string, create FactoryFunc) {
	if _, ok := r.modules[name]; ok {
		panic(fmt.Sprintf("module %q already registered", name))
	}
	r.modules[name] = create
}

func (r *Registry) Create(name string, cfg []byte) (mirror.Module, error) {
	c, ok := r.modules[name]
	if !ok {
		return nil, fmt.Errorf("module %q does not exist", name)
	}

	m, err := c(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating module %q: %s", name, err)
	}

	return m, nil
}
