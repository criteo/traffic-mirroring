package control

import (
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
)

type Virtual struct {
	ctx *mirror.ModuleContext
}

func NewVirtual(ctx *mirror.ModuleContext) mirror.Module {
	return &Virtual{ctx: ctx}
}

func (m *Virtual) Context() *mirror.ModuleContext {
	return m.ctx
}

func (m *Virtual) Children() [][]mirror.Module {
	return nil
}

func (m *Virtual) Output() <-chan mirror.Request {
	return nil
}

func (m *Virtual) SetInput(c <-chan mirror.Request) {

}
