package config

import (
	"io"
	"io/ioutil"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/modules/control"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
)

func Create(r io.Reader) (mirror.Module, error) {
	cfg, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return registry.Create(control.SeqName, cfg)
}
