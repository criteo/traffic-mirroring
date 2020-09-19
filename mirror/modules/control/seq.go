package control

import (
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/config"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

const (
	SeqName = "control.seq"
)

func init() {
	registry.Register(SeqName, NewSeq)
}

type Seq struct {
	ctx     *mirror.ModuleContext
	out     <-chan mirror.Request
	modules []mirror.Module
}

func NewSeq(ctx *mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	mods, err := config.CreateModules(cfg)
	if err != nil {
		return nil, err
	}

	var lastOut <-chan mirror.Request
	for _, sub := range mods {
		if lastOut != nil {
			sub.SetInput(lastOut)
		}

		lastOut = sub.Output()
	}

	return &Seq{
		ctx:     ctx,
		out:     lastOut,
		modules: mods,
	}, nil
}

func (m *Seq) Context() *mirror.ModuleContext {
	return m.ctx
}

func (m *Seq) Children() [][]mirror.Module {
	return [][]mirror.Module{m.modules}
}

func (m *Seq) Output() <-chan mirror.Request {
	return m.out
}

func (m *Seq) SetInput(c <-chan mirror.Request) {
	if len(m.modules) == 0 {
		return
	}

	m.modules[0].SetInput(c)
}
