package control

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/expr"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/modules"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

const (
	RateLimitName = "control.rate_limit"
)

func init() {
	registry.Register(RateLimitName, NewRateLimit)
}

type RateLimitConfig struct {
	RPS *expr.NumberExpr `json:"rps"`
}

type RateLimit struct {
	ctx mirror.ModuleContext
	cfg RateLimitConfig
	out chan mirror.Request

	ready    bool
	interval time.Duration
}

func NewRateLimit(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	c := RateLimitConfig{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	mod := &RateLimit{
		ctx: ctx,
		cfg: c,
		out: make(chan mirror.Request),
	}

	return mod, nil
}

func (m *RateLimit) Output() <-chan mirror.Request {
	return m.out
}

func (m *RateLimit) SetInput(c <-chan mirror.Request) {
	go func() {
		last := time.Now()
		for r := range c {
			if !m.ready {
				err := m.init(r)
				if err != nil {
					log.Fatalf("%s: %s", RateLimitName, err)
				}
			}

			now := time.Now()
			time.Sleep(m.interval - now.Sub(last))
			last = now

			modules.RequestsTotal.WithLabelValues(m.ctx.Name).Inc()
			m.out <- r
		}
		close(m.out)
	}()
}

func (m *RateLimit) init(r mirror.Request) error {
	v, err := m.cfg.RPS.EvalFloat(r)
	if err != nil {
		return fmt.Errorf("cannot evaluate rps: %s", err)
	}

	m.ready = true
	m.interval = time.Duration((1 / v) * float64(time.Second))
	return nil
}
