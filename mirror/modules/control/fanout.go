package control

import (
	"encoding/json"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

const (
	FanoutName = "control.fanout"
)

func init() {
	registry.Register(FanoutName, NewFanout)
}

type FanoutConfigEl struct {
	Type   string `json:"type"`
	Config json.RawMessage
}

type Fanout struct {
	out     chan mirror.Request
	modules []mirror.Module
	in      []chan mirror.Request
}

func NewFanout(cfg []byte) (mirror.Module, error) {
	mod := &Fanout{
		out: make(chan mirror.Request),
	}

	c := []FanoutConfigEl{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	for _, m := range c {
		sub, err := registry.Create(m.Type, []byte(m.Config))
		if err != nil {
			return nil, err
		}

		in := make(chan mirror.Request)
		mod.in = append(mod.in, in)
		sub.SetInput(in)
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
	}()
}
