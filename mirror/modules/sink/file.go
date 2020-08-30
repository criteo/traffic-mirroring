package sink

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gogo/protobuf/proto"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
	log "github.com/sirupsen/logrus"
)

const (
	FileName = "sink.file"
)

func init() {
	registry.Register(FileName, NewFile)
}

type FileConfig struct {
	Path   string `json:"path,omitempty"`
	Format string `json:"format,omitempty"`
}

type File struct {
	f   *bufio.Writer
	fmt string
	out chan mirror.Request

	encoder func(w io.Writer, req mirror.Request) error

	protoBuf     *proto.Buffer
	protoSizeBuf []byte
	jsonEnc      *json.Encoder
}

func NewFile(cfg []byte) (mirror.Module, error) {
	c := FileConfig{}

	err := json.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(c.Path, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriter(f)

	mod := &File{
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
