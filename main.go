package main

import (
	"flag"
	"os"

	"github.com/criteo/traffic-mirroring/mirror/config"
	_ "github.com/criteo/traffic-mirroring/mirror/modules/control"
	_ "github.com/criteo/traffic-mirroring/mirror/modules/sink"
	_ "github.com/criteo/traffic-mirroring/mirror/modules/source"
	"github.com/criteo/traffic-mirroring/mirror/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfgPath := flag.String("c", "config.json", "Config file path")
	logLevel := flag.String("log-level", "info", "Log level")
	flag.Parse()

	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(l)

	f, err := os.Open(*cfgPath)
	if err != nil {
		log.Fatalf("cannot open config file: %s", err)
	}

	cfg, err := config.Create(f)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ListenAddr != "" {
		srv := server.New(cfg.ListenAddr, cfg.Pipeline)
		go func() {
			err := srv.Run()
			if err != nil {
				log.Error(err)
			}
		}()
	}

	for range cfg.Pipeline.Output() {
	}
}
