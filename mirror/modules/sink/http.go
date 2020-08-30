package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/common/log"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

const (
	HTTPName = "sink.http"
)

func init() {
	registry.Register(HTTPName, NewHTTP)
}

type HTTPConfig struct {
	TargetURL string `json:"target_url,omitempty"`
	Timeout   string `json:"timeout,omitempty"`
}

type HTTP struct {
	out    chan mirror.Request
	client *http.Client
	target *url.URL
}

func NewHTTP(cfg []byte) (mirror.Module, error) {
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

	mod := &HTTP{
		out: make(chan mirror.Request),
		client: &http.Client{
			Timeout: timeout,
		},
		target: url,
	}
	return mod, nil
}

func (m *HTTP) Output() <-chan mirror.Request {
	return m.out
}

func (m *HTTP) SetInput(c <-chan mirror.Request) {
	go func() {
		for r := range c {
			m.sendRequest(r)
		}
		close(m.out)
	}()
}

func (m *HTTP) sendRequest(req mirror.Request) {
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

	res, err := m.client.Do(hreq)
	if err != nil {
		log.Errorf("%s: %q: %s", HTTPName, m.target.Host, err)
	}

	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
}
