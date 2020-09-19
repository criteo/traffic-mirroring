package control

import (
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

const (
	IdentityName = "control.identity"
)

func init() {
	registry.Register(IdentityName, NewIdentity)
}

type Identity struct {
	ctx *mirror.ModuleContext
	out chan mirror.Request
}

func NewIdentity(ctx *mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	return &Identity{
		out: make(chan mirror.Request),
		ctx: ctx,
	}, nil
}

func (m *Identity) Context() *mirror.ModuleContext {
	return m.ctx
}

func (m *Identity) Children() [][]mirror.Module {
	return nil
}

func (m *Identity) Output() <-chan mirror.Request {
	return m.out
}

func (m *Identity) SetInput(c <-chan mirror.Request) {
	go func() {
		for r := range c {
			m.ctx.HandledRequest()
			m.out <- r
		}
		close(m.out)
	}()
}
