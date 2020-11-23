package control

import (
	"sync/atomic"

	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/criteo/traffic-mirroring/mirror/config"
	"github.com/criteo/traffic-mirroring/mirror/registry"
)

const (
	FanoutName = "control.fanout"
)

func init() {
	registry.Register(FanoutName, NewFanout)
}

type Fanout struct {
	ctx     *mirror.ModuleContext
	out     chan mirror.Request
	modules []mirror.Module
	in      []chan mirror.Request

	outClosed uint32
}

func NewFanout(ctx *mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	mod := &Fanout{
		ctx: ctx,
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

	mod.modules = mods

	return mod, nil
}

func (m *Fanout) Context() *mirror.ModuleContext {
	return m.ctx
}

func (m *Fanout) Children() [][]mirror.Module {
	res := [][]mirror.Module{}
	for _, m := range m.modules {
		res = append(res, []mirror.Module{m})
	}
	return res
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
