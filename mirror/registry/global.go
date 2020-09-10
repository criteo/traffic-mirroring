package registry

import "github.com/shimmerglass/http-mirror-pipeline/mirror"

var DefaultRegistry = New()

func Register(moduleType string, create FactoryFunc) {
	DefaultRegistry.Register(moduleType, create)
}

func Create(moduleType string, ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	return DefaultRegistry.Create(moduleType, ctx, cfg)
}
