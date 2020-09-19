package mirror

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requets_total",
		Help: "The total number of requets handled by the module",
	}, []string{"module"})
)

type ModuleContext struct {
	Type           string
	Name           string
	RPS            int
	requestCounter uint64
}

func (c *ModuleContext) Role() string {
	return strings.Split(c.Type, ".")[0]
}

func (c *ModuleContext) HandledRequest() {
	atomic.AddUint64(&c.requestCounter, 1)
	RequestsTotal.WithLabelValues(c.Name).Inc()
}

func (c *ModuleContext) Run() {
	for range time.Tick(time.Second) {
		c.RPS = int(atomic.SwapUint64(&c.requestCounter, 0))
	}
}

type Module interface {
	Context() *ModuleContext
	SetInput(<-chan Request)
	Output() <-chan Request
	Children() [][]Module
}
