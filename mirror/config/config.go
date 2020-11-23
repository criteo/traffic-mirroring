package config

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/criteo/traffic-mirroring/mirror/registry"
)

type pipeline struct {
	inner mirror.Module
}

func (p *pipeline) Context() *mirror.ModuleContext {
	return p.inner.Context()
}

func (p *pipeline) Children() [][]mirror.Module {
	return p.inner.Children()
}

func (p *pipeline) SetInput(c <-chan mirror.Request) {
	p.inner.SetInput(c)
}
func (p *pipeline) Output() <-chan mirror.Request {
	return p.inner.Output()
}
func (p *pipeline) UnmarshalJSON(b []byte) error {
	ctx := &mirror.ModuleContext{
		Name: "Pipeline",
		Type: "virtual.pipeline",
	}
	m, err := registry.Create("control.seq", ctx, b)
	if err != nil {
		return err
	}

	p.inner = m
	return nil
}

type Config struct {
	ListenAddr string `json:"listen_addr,omitempty"`

	Pipeline *pipeline `json:"pipeline,omitempty"`
}

func Create(r io.Reader) (Config, error) {
	cfg := Config{}
	err := json.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("error reading config: %w", err)
	}
	return cfg, nil
}
