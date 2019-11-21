package main

import (
	"log"
)

func init() {
	// Init config
	initConfig()
}

func main() {
	p, err := NewProxy(cfg.Type, cfg.Source, cfg.Destination)
	if err != nil {
		log.Fatalf("couldn't start a TCP proxy: %v", err)
	}

	log.Printf("[*] Starting proxy of type: %s, listening on: %s, destination: %s",
		cfg.Type, cfg.Source, cfg.Destination)
	err = p.StartProxy()
	if err != nil {
		log.Fatalf("error while running proxy: %v", err)
	}
}
