package main

import (
	"flag"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

type tlsConfig struct {
	Domain             string `koanf:"domain"`
	CertsPath          string `koanf:"certs_path"`
	CloudflareAPIToken string `koanf:"cloudflare_api_token"`
	Email              string `koanf:"email"`
	DisableAutomatic   bool   `koanf:"disable_automatic"`
	Certificate        string `koanf:"cert"`
	Key                string `koanf:"key"`
}

type proxyConfig struct {
	Type        string `koanf:"type"`
	Source      string `koanf:"source"`
	Destination string `koanf:"destination"`
}

// Config has all the app config
type Config struct {
	Proxy proxyConfig
	TLS   tlsConfig
}

func initConfig() (Config, error) {
	var cfg Config
	var k = koanf.New(".")
	configPath := flag.String("config", "config.toml", "Path to configuration")
	flag.Parse()

	// Load TOML config.
	if err := k.Load(file.Provider(*configPath), toml.Parser()); err != nil {
		return Config{}, err
	}

	k.Unmarshal("proxy", &cfg.Proxy)
	k.Unmarshal("tls", &cfg.TLS)
	return cfg, nil
}
