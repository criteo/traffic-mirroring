package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/log"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/modules"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
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
	TargetURL string `json:"target_url,omitempty"`
	Timeout   string `json:"timeout,omitempty"`
	Parallel  int    `json:"parallel"`
}

type HTTP struct {
	ctx    mirror.ModuleContext
	out    chan mirror.Request
	tasks  chan mirror.Request
	client *http.Client
	target *url.URL

	numWorkers int
	maxWorkers int
}

func NewHTTP(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	c := HTTPConfig{}

	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("timeout: %w", err)
	}

	url, err := url.Parse(c.TargetURL)
	if err != nil {
		return nil, err
	}

	maxWorkers := 1
	if c.Parallel > 0 {
		maxWorkers = c.Parallel
	}

	mod := &HTTP{
		ctx:   ctx,
		out:   make(chan mirror.Request),
		tasks: make(chan mirror.Request),
		client: &http.Client{
			Timeout: timeout,
		},
		target:     url,
		maxWorkers: maxWorkers,
	}
	return mod, nil
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
				}
			}
		}
		close(m.out)
	}()
}

func (m *HTTP) sendRequest(req mirror.Request) {
	modules.RequestsTotal.WithLabelValues(m.ctx.Name).Inc()

	headers := http.Header{}
	for name, vals := range req.Headers {
		headers[name] = vals.Values
	}

	hreq := &http.Request{
		Method: req.Method.String(),
		URL: &url.URL{
			Scheme: m.target.Scheme,
			Host:   m.target.Host,
			Path:   req.Path,
		},
		Header: headers,
		Body:   ioutil.NopCloser(bytes.NewBuffer(req.Body)),
	}

	start := time.Now()
	res, err := m.client.Do(hreq)
	if err != nil {
		log.Errorf("%s: %q: %s", HTTPName, m.target.Host, err)
		return
	}

	httpResponseTime.WithLabelValues(m.ctx.Name).Observe(time.Since(start).Seconds())
	httpResponseTotal.WithLabelValues(m.ctx.Name, strconv.Itoa(res.StatusCode)).Inc()

	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
}

func (m *HTTP) runWorker(req mirror.Request) {
	m.sendRequest(req)
	timeout := time.NewTimer(workerTimeout)

	for {
		select {
		case req := <-m.tasks:
			m.sendRequest(req)
			timeout.Reset(workerTimeout)
		case <-timeout.C:
			m.numWorkers--
			break
		}
	}
}
