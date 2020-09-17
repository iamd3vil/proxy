package main

import (
	"fmt"
)

// Proxy needs to be implemented by any proxy
type Proxy interface {
	StartProxy() error
}

// NewProxy returns a new proxy instance according to proxy type
func NewProxy(config Config) (Proxy, error) {
	switch config.Proxy.Type {
	case "tcp":
		return NewTCPProxy(config.Proxy.Source, config.Proxy.Destination)
	case "tls":
		tc, err := getTLSConfig(config)
		if err != nil {
			return nil, err
		}
		return NewTLSProxy(config.Proxy.Source, config.Proxy.Destination, tc)
	default:
		return nil, fmt.Errorf("unrecognized proxy type: %s", config.Proxy.Type)
	}
}
