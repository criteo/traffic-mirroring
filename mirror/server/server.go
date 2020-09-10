package server

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	listenAddr string
}

func New(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/health", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	}))

	log.Infof("%s: listening on %s", os.Args[0], s.listenAddr)
	return http.ListenAndServe(s.listenAddr, mux)
}
