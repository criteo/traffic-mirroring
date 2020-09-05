package main

import (
	"flag"
	"os"

	"github.com/shimmerglass/http-mirror-pipeline/mirror/config"
	_ "github.com/shimmerglass/http-mirror-pipeline/mirror/modules/control"
	_ "github.com/shimmerglass/http-mirror-pipeline/mirror/modules/sink"
	_ "github.com/shimmerglass/http-mirror-pipeline/mirror/modules/source"
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

	module, err := config.Create(f)
	if err != nil {
		log.Fatal(err)
	}

	for range module.Output() {
	}
}
