package main

import (
	"log"
)

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("error while reading config: %v", err)
	}

	p, err := NewProxy(cfg)
	if err != nil {
		log.Fatalf("couldn't start a TCP proxy: %v", err)
	}

	log.Printf("[*] Starting proxy of type: %s, listening on: %s, destination: %s",
		cfg.Proxy.Type, cfg.Proxy.Source, cfg.Proxy.Destination)
	err = p.StartProxy()
	if err != nil {
		log.Fatalf("error while running proxy: %v", err)
	}
}
