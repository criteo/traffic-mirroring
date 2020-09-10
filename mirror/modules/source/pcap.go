package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/modules"
	"github.com/shimmerglass/http-mirror-pipeline/mirror/registry"
	log "github.com/sirupsen/logrus"
)

const (
	PCapName = "source.pcap"
)

func init() {
	registry.Register(PCapName, NewPCap)
}

type blockingReader struct {
	buf  bytes.Buffer
	cond *sync.Cond
}

func newBlockingReader() *blockingReader {
	m := sync.Mutex{}
	return &blockingReader{
		cond: sync.NewCond(&m),
		buf:  bytes.Buffer{},
	}
}

func (br *blockingReader) Write(b []byte) (ln int, err error) {
	ln, err = br.buf.Write(b)
	br.cond.Broadcast()
	return
}

func (br *blockingReader) Read(b []byte) (ln int, err error) {
	ln, err = br.buf.Read(b)
	if err == io.EOF {
		br.cond.L.Lock()
		br.cond.Wait()
		br.cond.L.Unlock()
		ln, err = br.buf.Read(b)
	}
	return
}

type PCapConfig struct {
	Interface string `json:"interface,omitempty"`
	Port      int    `json:"port,omitempty"`
}

type PCap struct {
	cfg  PCapConfig
	ctx  mirror.ModuleContext
	out  chan mirror.Request
	buff *bytes.Buffer
}

func NewPCap(ctx mirror.ModuleContext, cfg []byte) (mirror.Module, error) {
	mod := &PCap{
		ctx:  ctx,
		out:  make(chan mirror.Request),
		buff: &bytes.Buffer{},
	}

	err := json.Unmarshal(cfg, &mod.cfg)
	if err != nil {
		return nil, err
	}

	if len(mod.cfg.Interface) == 0 {
		return nil, errors.New("interface is required")
	}

	err = mod.start()
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func (m *PCap) Output() <-chan mirror.Request {
	return m.out
}

func (m *PCap) SetInput(c <-chan mirror.Request) {
	log.Fatalf("%s: connot accept input", PCapName)
}

func (m *PCap) start() error {
	handle, err := pcap.OpenLive(m.cfg.Interface, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}

	filter := fmt.Sprintf("tcp and dst port %d", m.cfg.Port)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		return err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	log.Infof("%s: capturing on %q with filter %q", PCapName, m.cfg.Interface, filter)

	reader := newBlockingReader()
	bufReader := bufio.NewReader(reader)

	go func() {
		for packet := range packetSource.Packets() {
			appLayer := packet.ApplicationLayer()
			if appLayer == nil {
				continue
			}

			reader.Write(appLayer.Payload())
		}
	}()

	go func() {
		for {
			req, err := m.readRequest(bufReader)
			if err != nil {
				log.Errorf("%s: %s", PCapName, err)
				continue
			}

			modules.RequestsTotal.WithLabelValues(m.ctx.Name).Inc()
			m.out <- req
		}
	}()

	return nil
}

func (m *PCap) readRequest(reader *bufio.Reader) (mirror.Request, error) {
	req, err := http.ReadRequest(reader)
	if err != nil {
		return mirror.Request{}, err
	}

	body, _ := ioutil.ReadAll(req.Body)

	mreq := mirror.Request{
		Time:    time.Now(),
		Method:  mirror.Method(mirror.Method_value[req.Method]),
		Path:    req.URL.Path,
		Body:    body,
		Headers: map[string]mirror.HeaderValues{},
	}

	for name, values := range req.Header {
		mreq.Headers[name] = mirror.HeaderValues{Values: values}
	}

	return mreq, nil
}
