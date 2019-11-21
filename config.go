package main

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

// Config has all the app config
type Config struct {
	Type        string `koanf:"type"`
	Source      string `koanf:"source"`
	Destination string `koanf:"destination"`
	Certificate string `koanf:"cert"`
	Key         string `koanf:"key"`
}

var cfg Config
var k = koanf.New(".")

func initConfig() {
	// Load TOML config.
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	k.Unmarshal("proxy", &cfg)
}
