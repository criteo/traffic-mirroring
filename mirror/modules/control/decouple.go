package control

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
	log "github.com/sirupsen/logrus"
)

const (
	DecoupleName       = "control.decouple"
	logDroppedInterval = 10 * time.Second
)

func init() {
	registry.Register(DecoupleName, NewDecouple)
}

type DecoupleConfig struct {
	Quiet bool `json:"log"`
}

type Decouple struct {
	out   chan mirror.Request
	quiet bool
}

func NewDecouple(cfg []byte) (mirror.Module, error) {
	c := DecoupleConfig{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	mod := &Decouple{
		out:   make(chan mirror.Request),
		quiet: c.Quiet,
	}

	return mod, nil
}

func (m *Decouple) Output() <-chan mirror.Request {
	return m.out
}

func (m *Decouple) SetInput(c <-chan mirror.Request) {
	dropped := uint32(0)

	go func() {
		for r := range c {
			select {
			case m.out <- r:
			default:
				atomic.AddUint32(&dropped, 1)
			}
		}
		close(m.out)
	}()

	if !m.quiet {
		go func() {
			for range time.Tick(logDroppedInterval) {
				v := atomic.SwapUint32(&dropped, 0)
				if v == 0 {
					continue
				}

				log.Warnf("%s: dropped %d requests", DecoupleName, v)
			}
		}()
	}
}
