package control

import (
	"encoding/json"
	"time"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
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
	RPS int `json:"rps"`
}

type RateLimit struct {
	ctx      mirror.ModuleContext
	out      chan mirror.Request
	interval time.Duration
}

func NewRateLimit(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	c := RateLimitConfig{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	mod := &RateLimit{
		ctx:      ctx,
		out:      make(chan mirror.Request),
		interval: time.Duration((1 / float64(c.RPS)) * float64(time.Second)),
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
			now := time.Now()
			time.Sleep(m.interval - now.Sub(last))
			last = now

			modules.RequestsTotal.WithLabelValues(m.ctx.Name).Inc()
			m.out <- r
		}
		close(m.out)
	}()
}
