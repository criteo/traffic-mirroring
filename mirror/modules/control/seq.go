package control

import (
	"encoding/json"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

const (
	SeqName = "control.seq"
)

func init() {
	registry.Register(SeqName, NewSeq)
}

type SeqConfigEl struct {
	Type   string `json:"type"`
	Config json.RawMessage
}

type Seq struct {
	out     <-chan mirror.Request
	modules []mirror.Module
}

func NewSeq(cfg []byte) (mirror.Module, error) {
	c := []SeqConfigEl{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	var lastOut <-chan mirror.Request
	for _, m := range c {
		sub, err := registry.Create(m.Type, []byte(m.Config))
		if err != nil {
			return nil, err
		}

		if lastOut != nil {
			sub.SetInput(lastOut)
		}

		lastOut = sub.Output()
	}

	return &Seq{out: lastOut}, nil
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
