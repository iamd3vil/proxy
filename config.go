package main

import (
	"flag"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

type tlsConfig struct {
	Domain             string        `koanf:"domain"`
	CertsPath          string        `koanf:"certs_path"`
	CloudflareAPIToken string        `koanf:"cloudflare_api_token"`
	Email              string        `koanf:"email"`
	DisableAutomatic   bool          `koanf:"disable_automatic"`
	Certificate        string        `koanf:"cert"`
	Key                string        `koanf:"key"`
	CheckForExpiry     time.Duration `koanf:"check_for_expiry"`
}

type proxyConfig struct {
	Type        string `koanf:"type"`
	Source      string `koanf:"source"`
	Destination string `koanf:"destination"`
}

// Config has all the app config
type Config struct {
	Proxy proxyConfig `koanf:"proxy"`
	TLS   tlsConfig   `koanf:"tls"`
}

func initConfig() (Config, error) {
	var cfg Config
	var k = koanf.New(".")
	configPath := flag.String("config", "config.toml", "Path to configuration")
	flag.Parse()

	k.Load(confmap.Provider(map[string]interface{}{
		"tls.check_for_expiry": 10 * 24 * time.Hour,
	}, "."), nil)

	// Load TOML config.
	if err := k.Load(file.Provider(*configPath), toml.Parser()); err != nil {
		return Config{}, err
	}

	k.Unmarshal("", &cfg)
	return cfg, nil
}
