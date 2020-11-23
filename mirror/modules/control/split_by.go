package control

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/criteo/traffic-mirroring/mirror/config"
	"github.com/criteo/traffic-mirroring/mirror/expr"
	"github.com/criteo/traffic-mirroring/mirror/registry"
	log "github.com/sirupsen/logrus"
)

const (
	SplitByName = "control.split_by"
)

func init() {
	registry.Register(SplitByName, NewSplitBy)
}

type SplitByConfig struct {
	Expr     *expr.AnyExpr   `json:"expr,omitempty"`
	Pipeline json.RawMessage `json:"pipeline,omitempty"`
}

type splitByModule struct {
	e             interface{}
	in            chan mirror.Request
	mod           mirror.Module
	lastMessageAt time.Time
}

type SplitBy struct {
	cfg         SplitByConfig
	ctx         *mirror.ModuleContext
	out         chan mirror.Request
	modules     map[interface{}]*splitByModule
	modulesLock sync.Mutex
}

func NewSplitBy(ctx *mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	c := SplitByConfig{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	// check that the module compiles
	_, err = config.CreateModule(c.Pipeline)
	if err != nil {
		return nil, err
	}

	mod := &SplitBy{
		cfg:     c,
		out:     make(chan mirror.Request),
		ctx:     ctx,
		modules: map[interface{}]*splitByModule{},
	}
	go mod.reapInactive()
	return mod, nil
}

func (m *SplitBy) Context() *mirror.ModuleContext {
	return m.ctx
}

func (m *SplitBy) Children() [][]mirror.Module {
	m.modulesLock.Lock()
	defer m.modulesLock.Unlock()

	res := [][]mirror.Module{}
	for _, m := range m.modules {
		res = append(res, []mirror.Module{
			NewVirtual(&mirror.ModuleContext{
				Type: "virtual.split_by",
				Name: fmt.Sprint(m.e),
			}),
			m.mod,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		a := res[i][0].Context().Name
		b := res[j][0].Context().Name

		return a < b
	})

	if len(res) == 0 {
		res = append(res, []mirror.Module{
			NewVirtual(&mirror.ModuleContext{
				Type: "virtual.split_by",
				Name: "Awaiting data",
			}),
		})
	}

	return res
}

func (m *SplitBy) Output() <-chan mirror.Request {
	return m.out
}

func (m *SplitBy) SetInput(c <-chan mirror.Request) {
	go func() {
		for r := range c {
			m.ctx.HandledRequest()

			m.modulesLock.Lock()
			in, err := m.selectIn(r)
			if err != nil {
				log.Errorf("%s: %s", SplitByName, err)
				m.modulesLock.Unlock()
				continue
			}

			in <- r
			m.modulesLock.Unlock()
		}
		close(m.out)
	}()
}

func (m *SplitBy) selectIn(r mirror.Request) (chan mirror.Request, error) {
	e, err := m.cfg.Expr.Eval(r)
	if err != nil {
		return nil, err
	}

	if mod, ok := m.modules[e]; ok {
		mod.lastMessageAt = time.Now()
		return mod.in, nil
	}

	log.Debugf("%s: creating pipeline for value %q", SplitByName, e)
	mod, err := config.CreateModule(m.cfg.Pipeline)
	if err != nil {
		return nil, err
	}
	go func() {
		for r := range mod.Output() {
			m.out <- r
		}
	}()

	smod := &splitByModule{
		e:             e,
		in:            make(chan mirror.Request),
		mod:           mod,
		lastMessageAt: time.Now(),
	}
	mod.SetInput(smod.in)
	m.modules[e] = smod

	return smod.in, nil
}

func (m *SplitBy) reapInactive() {
	for range time.Tick(time.Minute) {
		m.modulesLock.Lock()
		for k, mod := range m.modules {
			if time.Since(mod.lastMessageAt) > time.Minute {
				log.Debugf("%s: destroying incative pipeline for value %q", SplitByName, k)
				close(mod.in)
				delete(m.modules, k)
			}
		}
		m.modulesLock.Unlock()
	}
}
