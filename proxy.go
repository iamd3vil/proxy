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
			Certificates: []tls.Certificate{cer},
		}
		return NewTLSProxy(config.Source, config.Destination, tc)
	default:
		return nil, fmt.Errorf("unrecognized proxy type: %s", config.Type)
	}
}
