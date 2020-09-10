package sink

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gogo/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/expr"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/modules"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
	log "github.com/sirupsen/logrus"
)

const (
	FileName = "sink.file"
)

var (
	writtenTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "file_written_bytes_total",
		Help: "The total number of bytes written by the module",
	}, []string{"module"})
)

func init() {
	registry.Register(FileName, NewFile)
}

type metricsWriter struct {
	name  string
	inner io.Writer
}

func (m metricsWriter) Write(b []byte) (int, error) {
	n, err := m.inner.Write(b)
	writtenTotal.WithLabelValues(m.name).Add(float64(n))
	return n, err
}

type FileConfig struct {
	Path       *expr.StringExpr `json:"path,omitempty"`
	Format     string           `json:"format,omitempty"`
	BufferSize int              `json:"buffer_size,omitempty"`
}

type File struct {
	ctx mirror.ModuleContext
	cfg FileConfig

	f   *bufio.Writer
	fmt string
	out chan mirror.Request

	encoder func(w io.Writer, req mirror.Request) error

	protoBuf     *proto.Buffer
	protoSizeBuf []byte
	jsonEnc      *json.Encoder

	ready bool
}

func NewFile(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	c := FileConfig{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	mod := &File{
		ctx: ctx,
		cfg: c,
		out: make(chan mirror.Request),
	}

	if c.Path.Static() {
		err := mod.init(mirror.Request{})
		if err != nil {
			return nil, err
		}
	}

	return mod, nil
}

func (m *File) Output() <-chan mirror.Request {
	return m.out
}

func (m *File) SetInput(c <-chan mirror.Request) {
	go func() {
		for r := range c {
			if !m.ready {
				err := m.init(r)
				if err != nil {
					log.Fatalf("%s: %s", FileName, err)
				}
			}

			modules.RequestsTotal.WithLabelValues(m.ctx.Name).Inc()
			err := m.encoder(m.f, r)
			if err != nil {
				log.Errorf("%s: %s", FileName, err)
			}
		}
		err := m.f.Flush()
		if err != nil {
			log.Errorf("%s: %s", FileName, err)
		}
		close(m.out)
	}()
}

func (m *File) writeJSON(w io.Writer, req mirror.Request) error {
	err := m.jsonEnc.Encode(req)
	if err != nil {
		return err
	}

	return nil
}

func (m *File) writeProto(w io.Writer, req mirror.Request) error {
	m.protoBuf.Reset()
	err := m.protoBuf.Marshal(&req)
	if err != nil {
		return err
	}

	n := binary.PutUvarint(m.protoSizeBuf, uint64(len(m.protoBuf.Bytes())))

	_, err = w.Write(m.protoSizeBuf[:n])
	if err != nil {
		return err
	}

	_, err = w.Write(m.protoBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (m *File) init(r mirror.Request) error {
	path, err := m.cfg.Path.Eval(r)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	if m.cfg.BufferSize == 0 {
		m.cfg.BufferSize = 1024
	}

	w := bufio.NewWriterSize(metricsWriter{
		name:  m.ctx.Name,
		inner: f,
	}, m.cfg.BufferSize)

	m.f = w
	m.protoBuf = proto.NewBuffer(nil)
	m.protoSizeBuf = make([]byte, binary.MaxVarintLen64)
	m.jsonEnc = json.NewEncoder(w)

	switch m.cfg.Format {
	case "proto":
		m.encoder = m.writeProto
	case "json":
		m.encoder = m.writeJSON
	default:
		return fmt.Errorf("unknown format %q", m.cfg.Format)
	}

	m.ready = true

	return nil
}
