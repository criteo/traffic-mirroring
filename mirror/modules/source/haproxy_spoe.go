package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	spoe "github.com/criteo/haproxy-spoe-go"
	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/criteo/traffic-mirroring/mirror/registry"
	log "github.com/sirupsen/logrus"
)

const (
	HAProxySPOEName        = "source.haproxy_spoe"
	DefaultSpoeIdleTimeout = 30
)

func init() {
	registry.Register(HAProxySPOEName, NewHAProxySPOE)
}

type mappingFunc func(req *mirror.Request, value interface{}) error

var defaultMappingConfig = map[string]string{
	"method":  "method",
	"ver":     "ver",
	"path":    "path",
	"headers": "headers",
	"body":    "body",
}

type HAProxySPOEConfig struct {
	ListenAddr  string `json:"listen_addr"`
	IdleTimeout string `json:"idle_timeout"`
	Mapping     map[string]string
}

type HAProxySPOE struct {
	cfg         HAProxySPOEConfig
	ctx         *mirror.ModuleContext
	out         chan mirror.Request
	mapping     map[string]mappingFunc
	idleTimeout time.Duration
}

func NewHAProxySPOE(ctx *mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	mod := &HAProxySPOE{
		ctx:     ctx,
		out:     make(chan mirror.Request),
		mapping: map[string]mappingFunc{},
	}

	err := json.Unmarshal(cfg, &mod.cfg)
	if err != nil {
		return nil, err
	}

	if len(mod.cfg.ListenAddr) == 0 {
		return nil, errors.New("listen_addr is required")
	}

	if mod.cfg.IdleTimeout != "" {
		t, err := time.ParseDuration(mod.cfg.IdleTimeout)
		if err != nil {
			return nil, errors.New("idle_timeout parse")
		}
		mod.idleTimeout = t * time.Second
	} else {
		mod.idleTimeout = DefaultSpoeIdleTimeout * time.Second
	}

	mappingCfg := defaultMappingConfig
	for k, v := range mod.cfg.Mapping {
		mappingCfg[k] = v
	}

	for k, v := range mappingCfg {
		switch v {
		case "method":
			mod.mapping[k] = mapMethod
		case "ver":
			mod.mapping[k] = mapVer
		case "path":
			mod.mapping[k] = mapPath
		case "headers":
			mod.mapping[k] = mapHeaders
		case "body":
			mod.mapping[k] = mapBody
		default:
			if !strings.HasPrefix(v, "meta.") {
				return nil, fmt.Errorf("unknown mapping key %q", v)
			}

			mod.mapping[k] = mapMeta(strings.TrimPrefix(v, "meta."))
		}
	}

	err = mod.start()
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func (m *HAProxySPOE) Context() *mirror.ModuleContext {
	return m.ctx
}

func (m *HAProxySPOE) Children() [][]mirror.Module {
	return nil
}

func (m *HAProxySPOE) Output() <-chan mirror.Request {
	return m.out
}

func (m *HAProxySPOE) SetInput(c <-chan mirror.Request) {
	log.Fatalf("%s: connot accept input", HAProxySPOEName)
}

func (m *HAProxySPOE) start() error {
	agent := spoe.NewWithConfig(m.handleMessage, spoe.Config{
		IdleTimeout: m.idleTimeout,
	})

	var l net.Listener
	var err error
	if m.cfg.ListenAddr[0] == '@' {
		syscall.Unlink(m.cfg.ListenAddr[1:])
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

func (m *HAProxySPOE) handleMessage(msgs *spoe.MessageIterator) ([]spoe.Action, error) {
	for msgs.Next() {
		msg := msgs.Message

		req := mirror.Request{}

		for msg.Args.Next() {
			arg := msg.Args.Arg

			m, ok := m.mapping[arg.Name]
			if !ok {
				continue
			}

			err := m(&req, arg.Value)
			if err != nil {
				log.Errorf("%s: bad message: %s", HAProxySPOEName, err)
			}
		}

		m.ctx.HandledRequest()
		m.out <- req
	}

	if err := msgs.Error(); err != nil {
		log.Errorf("%s: error handling message: %s", HAProxySPOEName, err)
	}

	return nil, nil
}

func mapMethod(req *mirror.Request, value interface{}) error {
	method, ok := value.(string)
	if !ok {
		return fmt.Errorf("bad type %T received for method, expected string", value)
	}

	i, ok := mirror.Method_value[method]
	if !ok {
		return fmt.Errorf("unknown method %q", value)
	}

	req.Method = mirror.Method(i)
	return nil
}

func mapPath(req *mirror.Request, value interface{}) error {
	path, ok := value.(string)
	if !ok {
		return fmt.Errorf("bad type %T received for path, expected string", value)
	}

	req.Path = path
	return nil
}

func mapVer(req *mirror.Request, value interface{}) error {
	ver, ok := value.(string)
	if !ok {
		return fmt.Errorf("bad type %T received for http version, expected string", value)
	}

	switch ver {
	case "1.1":
		req.HttpVersion = mirror.HTTPVersion_HTTP1_1
	case "2":
		req.HttpVersion = mirror.HTTPVersion_HTTP2
	default:
		req.HttpVersion = mirror.HTTPVersion_HTTP1_0
	}

	return nil
}

func mapHeaders(req *mirror.Request, value interface{}) error {
	rawHeaders, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("bad type %T received for headers, expected bytes", value)
	}

	headers, err := spoe.DecodeHeaders(rawHeaders)
	if err != nil {
		return err
	}

	req.Headers = map[string]*mirror.HeaderValue{}
	for name, values := range headers {
		req.Headers[name] = &mirror.HeaderValue{Values: values}
	}

	return nil
}

func mapBody(req *mirror.Request, value interface{}) error {
	body, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("bad type %T received for body, expected bytes", value)
	}

	req.Body = body

	return nil
}

func mapMeta(name string) mappingFunc {
	return func(req *mirror.Request, value interface{}) error {
		if req.Meta == nil {
			req.Meta = map[string]*mirror.MetaValue{}
		}

		switch t := value.(type) {
		case string:
			req.Meta[name] = &mirror.MetaValue{
				Value: &mirror.MetaValue_String_{t},
			}
		case int:
			req.Meta[name] = &mirror.MetaValue{
				Value: &mirror.MetaValue_Int{int64(t)},
			}
		case bool:
			req.Meta[name] = &mirror.MetaValue{
				Value: &mirror.MetaValue_Bool{t},
			}

		default:
			return fmt.Errorf("unhandled type for meta %T, allowed values are string, int, bool", value)
		}

		return nil
	}
}
