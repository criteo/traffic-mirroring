package control

import (
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/modules"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

const (
	IdentityName = "control.identity"
)

func init() {
	registry.Register(IdentityName, NewIdentity)
}

type Identity struct {
	ctx     mirror.ModuleContext
	out     chan mirror.Request
	modules []mirror.Module
}

func NewIdentity(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	return &Identity{
		out: make(chan mirror.Request),
	}, nil
}

func (m *Identity) Output() <-chan mirror.Request {
	return m.out
}

func (m *Identity) SetInput(c <-chan mirror.Request) {
	go func() {
		for r := range c {
			modules.RequestsTotal.WithLabelValues(m.ctx.Name).Inc()
			m.out <- r
		}
		close(m.out)
	}()
}
