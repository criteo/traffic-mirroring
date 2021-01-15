package control

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/criteo/traffic-mirroring/mirror/registry"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

const (
	DecoupleName       = "control.decouple"
	logDroppedInterval = 10 * time.Second
)

var (
	droppedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "decouple_dropped_total",
		Help: "The total number of responses dropped",
	}, []string{"module"})
)

func init() {
	registry.Register(DecoupleName, NewDecouple)
}

type DecoupleConfig struct {
	Quiet     bool `json:"log"`
	QueueSize int  `json:"queue_size"`
}

type Decouple struct {
	ctx   *mirror.ModuleContext
	out   chan mirror.Request
	quiet bool
}

func NewDecouple(ctx *mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	c := DecoupleConfig{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	queueSize := 100
	if c.QueueSize > 0 {
		queueSize = c.QueueSize
	}

	mod := &Decouple{
		ctx:   ctx,
		out:   make(chan mirror.Request, queueSize),
		quiet: c.Quiet,
	}

	return mod, nil
}

func (m *Decouple) Context() *mirror.ModuleContext {
	return m.ctx
}

func (m *Decouple) Children() [][]mirror.Module {
	return nil
}

func (m *Decouple) Output() <-chan mirror.Request {
	return m.out
}

func (m *Decouple) SetInput(c <-chan mirror.Request) {
	dropped := uint32(0)
	processed := uint32(0)

	go func() {
		for r := range c {
			atomic.AddUint32(&processed, 1)
			m.ctx.HandledRequest()
			select {
			case m.out <- r:
			default:
				droppedTotal.WithLabelValues(m.ctx.Name).Inc()
				atomic.AddUint32(&dropped, 1)
			}
		}
		close(m.out)
	}()

	if !m.quiet {
		go func() {
			for range time.Tick(logDroppedInterval) {
				d := atomic.SwapUint32(&dropped, 0)
				p := atomic.SwapUint32(&processed, 0)
				if d == 0 {
					log.Debugf("%s: dropped %d of %d requests", DecoupleName, d, p)
				} else {
					log.Warnf("%s: dropped %d of %d requests", DecoupleName, d, p)
				}

			}
		}()
	}
}
