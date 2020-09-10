package control

import (
	"sync/atomic"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/config"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

const (
	FanoutName = "control.fanout"
)

func init() {
	registry.Register(FanoutName, NewFanout)
}

type Fanout struct {
	out     chan mirror.Request
	modules []mirror.Module
	in      []chan mirror.Request

	outClosed uint32
}

func NewFanout(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	mod := &Fanout{
		out: make(chan mirror.Request),
	}

	mods, err := config.CreateModules(cfg)
	if err != nil {
		return nil, err
	}

	for _, sub := range mods {
		in := make(chan mirror.Request)
		mod.in = append(mod.in, in)
		sub.SetInput(in)
		go mod.consume(sub)
	}

	return mod, nil
}

func (m *Fanout) Output() <-chan mirror.Request {
	return m.out
}

func (m *Fanout) SetInput(c <-chan mirror.Request) {
	go func() {
		for r := range c {
			for _, i := range m.in {
				i <- r
			}
		}

		for _, i := range m.in {
			close(i)
		}
	}()
}

func (m *Fanout) consume(mod mirror.Module) {
	out := mod.Output()
	for req := range out {
		m.out <- req
	}

	if int(atomic.AddUint32(&m.outClosed, 1)) == len(m.in) {
		close(m.out)
	}
}
