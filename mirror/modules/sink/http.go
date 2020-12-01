package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/criteo/traffic-mirroring/mirror/registry"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/log"
	"github.com/criteo/traffic-mirroring/mirror/expr"
)

const (
	HTTPName      = "sink.http"
	workerTimeout = 2 * time.Second
)

var (
	httpResponseTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_response_total",
		Help: "The total number of responses by status code",
	}, []string{"module", "status_code"})

	httpResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Http response time",
		Buckets: prometheus.ExponentialBuckets(0.001, 4, 5),
	}, []string{"module"})
)

func init() {
	registry.Register(HTTPName, NewHTTP)
}

type HTTPConfig struct {
	TargetURL *expr.StringExpr `json:"target_url,omitempty"`
	Timeout   string           `json:"timeout,omitempty"`
	Parallel  int              `json:"parallel"`
}

type HTTP struct {
	ctx    *mirror.ModuleContext
	out    chan mirror.Request
	tasks  chan mirror.Request
	client *http.Client
	cfg    HTTPConfig

	numWorkers int
	maxWorkers int
	workersWG  sync.WaitGroup
}

func NewHTTP(ctx *mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	c := HTTPConfig{}

	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("timeout: %w", err)
	}

	maxWorkers := 1
	if c.Parallel > 0 {
		maxWorkers = c.Parallel
	}

	mod := &HTTP{
		ctx:   ctx,
		cfg:   c,
		out:   make(chan mirror.Request),
		tasks: make(chan mirror.Request),
		client: &http.Client{
			Timeout: timeout,
		},
		maxWorkers: maxWorkers,
	}
	return mod, nil
}

func (m *HTTP) Context() *mirror.ModuleContext {
	return m.ctx
}

func (m *HTTP) Children() [][]mirror.Module {
	return nil
}

func (m *HTTP) Output() <-chan mirror.Request {
	return m.out
}

func (m *HTTP) SetInput(c <-chan mirror.Request) {
	go func() {
		for r := range c {
			if m.numWorkers >= m.maxWorkers {
				m.tasks <- r
			} else {
				select {
				case m.tasks <- r:
				default:
					go m.runWorker(r)
					m.numWorkers++
					m.workersWG.Add(1)
				}
			}
		}
		close(m.tasks)
		m.workersWG.Wait()
		close(m.out)
	}()
}

func (m *HTTP) sendRequest(req mirror.Request) {
	m.ctx.HandledRequest()
	baseURL, err := m.cfg.TargetURL.Eval(req)
	if err != nil {
		log.Errorf("%s: could not evaluate target URL: %s", HTTPName, err)
		return
	}
	url := baseURL + req.Path

	headers := http.Header{}
	for name, vals := range req.Headers {
		headers[name] = vals.Values
	}

	hreq, err := http.NewRequest(
		req.Method.String(),
		url,
		ioutil.NopCloser(bytes.NewBuffer(req.Body)),
	)
	if err != nil {
		log.Errorf("%s: could not create request: %s", HTTPName, err)
		return
	}
	hreq.Header = headers

	start := time.Now()
	res, err := m.client.Do(hreq)
	if err != nil {
		log.Errorf("%s: %q: %s", HTTPName, url, err)
		return
	}

	httpResponseTime.WithLabelValues(m.ctx.Name).Observe(time.Since(start).Seconds())
	httpResponseTotal.WithLabelValues(m.ctx.Name, strconv.Itoa(res.StatusCode)).Inc()

	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
}

func (m *HTTP) runWorker(req mirror.Request) {
	defer m.workersWG.Done()

	m.sendRequest(req)
	timeout := time.NewTimer(workerTimeout)

	for {
		select {
		case req, ok := <-m.tasks:
			if !ok {
				return
			}

			m.sendRequest(req)
			timeout.Reset(workerTimeout)
		case <-timeout.C:
			m.numWorkers--
			return
		}
	}
}
