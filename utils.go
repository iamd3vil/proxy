package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/libdns/cloudflare"
)

func getTLSConfig(config Config) (*tls.Config, error) {
	if !config.TLS.DisableAutomatic {
		magic := certmagic.NewDefault()
		magic.Storage = &certmagic.FileStorage{
			Path: config.TLS.CertsPath,
		}

		letsEncryptACME := certmagic.NewACMEIssuer(magic, certmagic.ACMEIssuer{
			CA:     certmagic.LetsEncryptProductionCA,
			Email:  config.TLS.Email,
			Agreed: true,
			DNS01Solver: &certmagic.DNS01Solver{
				DNSProvider: &cloudflare.Provider{
					APIToken: config.TLS.CloudflareAPIToken,
				},
			},
		})

		zerosslACME := certmagic.NewACMEIssuer(magic, certmagic.ACMEIssuer{
			CA:     certmagic.ZeroSSLProductionCA,
			Email:  config.TLS.Email,
			Agreed: true,
			DNS01Solver: &certmagic.DNS01Solver{
				DNSProvider: &cloudflare.Provider{
					APIToken: config.TLS.CloudflareAPIToken,
				},
			},
		})

		magic.Issuers = append(magic.Issuers, letsEncryptACME, zerosslACME)

		err := magic.ManageSync(context.Background(), []string{config.TLS.Domain})
		if err != nil {
			return nil, err
		}

		// Spin up a goroutine to check for cert expiry and renew if it expired.
		log.Printf("running a background process to renew certs every %s", config.TLS.CheckForExpiry)
		tc := time.NewTicker(config.TLS.CheckForExpiry)

		go func() {
			defer tc.Stop()
			for range tc.C {
				err := magic.ManageSync(context.Background(), []string{config.TLS.Domain})
				if err != nil {
					log.Printf("manage sync failed: failed to renew certs: %v", err)
				}
			}
		}()

		return magic.TLSConfig(), nil
	}
	// Read cert and key from file
	cert, err := os.ReadFile(config.TLS.Certificate)
	if err != nil {
		return nil, fmt.Errorf("error while reading cert: %v", err)
	}

	key, err := os.ReadFile(config.TLS.Key)
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
