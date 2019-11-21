package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
)

// Proxy needs to be implemented by any proxy
type Proxy interface {
	StartProxy() error
}

// NewProxy returns a new proxy instance according to proxy type
func NewProxy(config Config) (Proxy, error) {
	switch config.Type {
	case "tcp":
		return NewTCPProxy(config.Source, config.Destination)
	case "tls":
		// Read cert and key from file
		cert, err := ioutil.ReadFile(config.Certificate)
		if err != nil {
			return nil, fmt.Errorf("error while reading cert: %v", err)
		}

		key, err := ioutil.ReadFile(config.Key)
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
		return NewTLSProxy(config.Source, config.Destination, tc)
	default:
		return nil, fmt.Errorf("unrecognized proxy type: %s", config.Type)
	}
}
