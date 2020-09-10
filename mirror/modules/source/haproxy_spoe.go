package source

import (
	"encoding/json"
	"errors"
	"net"

	spoe "github.com/criteo/haproxy-spoe-go"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/modules"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
	log "github.com/sirupsen/logrus"
)

const (
	HAProxySPOEName = "source.haproxy_spoe"
)

func init() {
	registry.Register(HAProxySPOEName, NewHAProxySPOE)
}

type HAProxySPOEConfig struct {
	ListenAddr string `json:"listen_addr"`
}

type HAProxySPOE struct {
	cfg HAProxySPOEConfig
	ctx mirror.ModuleContext
	out chan mirror.Request
}

func NewHAProxySPOE(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	mod := &HAProxySPOE{
		ctx: ctx,
		out: make(chan mirror.Request),
	}

	err := json.Unmarshal(cfg, &mod.cfg)
	if err != nil {
		return nil, err
	}

	if len(mod.cfg.ListenAddr) == 0 {
		return nil, errors.New("listen_addr is required")
	}

	err = mod.start()
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func (m *HAProxySPOE) Output() <-chan mirror.Request {
	return m.out
}

func (m *HAProxySPOE) SetInput(c <-chan mirror.Request) {
	log.Fatalf("%s: connot accept input", HAProxySPOEName)
}

func (m *HAProxySPOE) start() error {
	agent := spoe.New(m.handleMessage)

	var l net.Listener
	var err error
	if m.cfg.ListenAddr[0] == '@' {
		l, err = net.Listen("unix", m.cfg.ListenAddr[1:])
	} else {
		l, err = net.Listen("tcp", m.cfg.ListenAddr)
	}
	if err != nil {
		return err
	}

	go func() {
		err := agent.Serve(l)
		if err != nil {
			log.Errorf("%s: %s", HAProxySPOEName, err)
		}
		close(m.out)
	}()

	return nil
}

func (m *HAProxySPOE) handleMessage(args []spoe.Message) ([]spoe.Action, error) {
	for _, msg := range args {
		method, ok := msg.Args["method"].(string)
		if !ok {
			logBadMessage(`missing "method"`)
			continue
		}
		path, ok := msg.Args["path"].(string)
		if !ok {
			logBadMessage(`missing "path"`)
			continue
		}
		ver, ok := msg.Args["ver"].(string)
		if !ok {
			logBadMessage(`missing "ver"`)
			continue
		}
		rawHeaders, ok := msg.Args["headers"].([]byte)
		if !ok {
			logBadMessage(`missing "headers"`)
			continue
		}
		headers, err := spoe.DecodeHeaders(rawHeaders)
		if err != nil {
			logBadMessage(err.Error())
			continue
		}
		body, ok := msg.Args["body"].([]byte)
		if !ok {
			logBadMessage(`missing "body"`)
			continue
		}

		req := mirror.Request{
			Method:  mirror.Method(mirror.Method_value[method]),
			Path:    path,
			Headers: map[string]mirror.HeaderValues{},
			Body:    body,
		}

		switch ver {
		case "1.1":
			req.HTTPVersion = mirror.V1_1
		case "2":
			req.HTTPVersion = mirror.V2
		default:
			req.HTTPVersion = mirror.V1_0
		}

		for name, values := range headers {
			req.Headers[name] = mirror.HeaderValues{Values: values}
		}

		modules.RequestsTotal.WithLabelValues(m.ctx.Name).Inc()
		m.out <- req
	}

	return nil, nil
}

func logBadMessage(msg string) {
	log.Errorf("%s: bad message: %s", HAProxySPOEName, msg)
}
