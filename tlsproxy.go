package main

import (
	"crypto/tls"
	"net"

	"inet.af/tcpproxy"
)

// TLSProxy represents a TLS Proxy
type TLSProxy struct {
	src, dst  string
	tlsConfig *tls.Config
}

// NewTLSProxy returns a new proxy instance
func NewTLSProxy(src, dst string, tlsConfig *tls.Config) (Proxy, error) {
	return &TLSProxy{
		src:       src,
		dst:       dst,
		tlsConfig: tlsConfig,
	}, nil
}

// StartProxy starts a TCP Proxy
func (p *TLSProxy) StartProxy() error {
	var proxy tcpproxy.Proxy
	proxy.ListenFunc = func(net, laddr string) (net.Listener, error) {
		return tls.Listen(net, laddr, p.tlsConfig)
	}
	proxy.AddRoute(p.src, tcpproxy.To(p.dst))
	return proxy.Run()
}
