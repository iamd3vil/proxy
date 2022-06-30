package main

import "inet.af/tcpproxy"

// TCPProxy represents a TCP Proxy
type TCPProxy struct {
	src, dst string
}

// NewTCPProxy returns a new proxy instance
func NewTCPProxy(src, dst string) (Proxy, error) {
	return &TCPProxy{
		src, dst,
	}, nil
}

// StartProxy starts a TCP Proxy
func (p *TCPProxy) StartProxy() error {
	var proxy tcpproxy.Proxy
	proxy.AddRoute(p.src, tcpproxy.To(p.dst))
	return proxy.Run()
}
