package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/caddyserver/certmagic"
	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
)

func getTLSConfig(config Config) (*tls.Config, error) {
	if !config.TLS.DisableAutomatic {
		// certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
		certmagic.Default.Storage = &certmagic.FileStorage{
			Path: config.TLS.CertsPath,
		}

		cfg := cloudflare.NewDefaultConfig()
		cfg.AuthToken = config.TLS.CloudflareAPIToken
		provider, err := cloudflare.NewDNSProviderConfig(cfg)
		if err != nil {
			return nil, err
		}

		certmagic.DefaultACME.Email = config.TLS.Email
		certmagic.DefaultACME.DNSProvider = provider
		// Get tls.Config from certmagic
		tlsConfig, err := certmagic.TLS([]string{config.TLS.Domain})
		if err != nil {
			return nil, err
		}
		return tlsConfig, nil
	}
	// Read cert and key from file
	cert, err := ioutil.ReadFile(config.TLS.Certificate)
	if err != nil {
		return nil, fmt.Errorf("error while reading cert: %v", err)
	}

	key, err := ioutil.ReadFile(config.TLS.Key)
	if err != nil {
		return nil, fmt.Errorf("error while reading key: %v", err)
	}

	cer, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("invalid certs: %v", err)
	}

	tc := &tls.Config{
		Certificates:             []tls.Certificate{cer},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		SessionTicketsDisabled: true,
	}
	return tc, nil
}
