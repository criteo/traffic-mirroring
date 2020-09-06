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
	out     chan mirror.Request
	modules []mirror.Module
}

func NewIdentity(cfg []byte) (mirror.Module, error) {
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
			m.out <- r
		}
		close(m.out)
	}()
}
