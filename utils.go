package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/caddyserver/certmagic"
	"github.com/libdns/cloudflare"
)

func getTLSConfig(config Config) (*tls.Config, error) {
	if !config.TLS.DisableAutomatic {
		magic := certmagic.NewDefault()
		magic.Storage = &certmagic.FileStorage{
			Path: config.TLS.CertsPath,
		}

		myACME := certmagic.NewACMEIssuer(magic, certmagic.ACMEIssuer{
			CA:     certmagic.LetsEncryptProductionCA,
			Email:  config.TLS.Email,
			Agreed: true,
			DNS01Solver: &certmagic.DNS01Solver{
				DNSProvider: &cloudflare.Provider{
					APIToken: config.TLS.CloudflareAPIToken,
				},
			},
		})

		magic.Issuers = append(magic.Issuers, myACME)

		err := magic.ManageSync(context.Background(), []string{config.TLS.Domain})
		if err != nil {
			return nil, err
		}

		return magic.TLSConfig(), nil
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
