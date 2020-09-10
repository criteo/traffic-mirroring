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
	Path       string `json:"path,omitempty"`
	Format     string `json:"format,omitempty"`
	BufferSize int    `json:"buffer_size,omitempty"`
}

type File struct {
	ctx mirror.ModuleContext
	f   *bufio.Writer
	fmt string
	out chan mirror.Request

	encoder func(w io.Writer, req mirror.Request) error

	protoBuf     *proto.Buffer
	protoSizeBuf []byte
	jsonEnc      *json.Encoder
}

func NewFile(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	c := FileConfig{}
	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(c.Path, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}

	if c.BufferSize == 0 {
		c.BufferSize = 1024
	}

	w := bufio.NewWriterSize(metricsWriter{
		name:  ctx.Name,
		inner: f,
	}, c.BufferSize)

	mod := &File{
		ctx: ctx,
		out: make(chan mirror.Request),
		fmt: c.Format,

		f:            w,
		protoBuf:     proto.NewBuffer(nil),
		protoSizeBuf: make([]byte, binary.MaxVarintLen64),
		jsonEnc:      json.NewEncoder(w),
	}

	switch c.Format {
	case "proto":
		mod.encoder = mod.writeProto
	case "json":
		mod.encoder = mod.writeJSON
	default:
		return nil, fmt.Errorf("unknown format %q", c.Format)
	}

	return mod, nil
}

func (m *File) Output() <-chan mirror.Request {
	return m.out
}

func (m *File) SetInput(c <-chan mirror.Request) {
	go func() {
		for r := range c {
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
