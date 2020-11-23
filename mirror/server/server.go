package server

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/criteo/traffic-mirroring/mirror/graph"
	_ "github.com/criteo/traffic-mirroring/mirror/server/statik"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	listenAddr string
	pipeline   mirror.Module
}

func New(listenAddr string, pipeline mirror.Module) *Server {
	return &Server{
		listenAddr: listenAddr,
		pipeline:   pipeline,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/health", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	}))

	mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/web/", http.StatusTemporaryRedirect)
	}))
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(statikFS)))
	mux.Handle("/api/graph", http.HandlerFunc(s.graphHandler))

	log.Infof("%s: listening on %s", os.Args[0], s.listenAddr)
	return http.ListenAndServe(s.listenAddr, mux)
}

func (s *Server) graphHandler(rw http.ResponseWriter, r *http.Request) {
	g := graph.FromModule(s.pipeline)
	rw.Header().Add("Content-Type", "text/plain")
	rw.Write([]byte(g.String()))
}
