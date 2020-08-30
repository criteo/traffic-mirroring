package registry

import "github.com/shimmerglass/http-mirror-pipeline/mirror"

var DefaultRegistry = New()

func Register(name string, create FactoryFunc) {
	DefaultRegistry.Register(name, create)
}

func Create(name string, cfg []byte) (mirror.Module, error) {
	return DefaultRegistry.Create(name, cfg)
}
