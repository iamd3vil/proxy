package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
)

// TLSProxy represents a TLS Proxy
type TLSProxy struct {
	Src       *net.TCPAddr
	Dst       *net.TCPAddr
	tlsConfig *tls.Config
}

// NewTLSProxy returns a new proxy instance
func NewTLSProxy(src, dst string, tlsConfig *tls.Config) (*TLSProxy, error) {
	laddr, err := net.ResolveTCPAddr("tcp", cfg.Source)
	if err != nil {
		return nil, err
	}

	raddr, err := net.ResolveTCPAddr("tcp", cfg.Destination)
	if err != nil {
		return nil, err
	}

	return &TLSProxy{
		Src:       laddr,
		Dst:       raddr,
		tlsConfig: tlsConfig,
	}, nil
}

// StartProxy starts a TCP Proxy
func (p *TLSProxy) StartProxy() error {
	listener, err := tls.Listen("tcp", p.Src.String(), p.tlsConfig)
	if err != nil {
		return fmt.Errorf("error while listening: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("couldn't accept connection: '%v'", err)
			continue
		}

		defer conn.Close()

		rconn, err := net.DialTCP("tcp", nil, p.Dst)
		if err != nil {
			conn.Close()
			log.Printf("couldn't connect to destination '%v': %v", p.Dst, err)
			continue
		}
		defer rconn.Close()

		closeChan := make(chan int, 2)

		go func() {
			io.Copy(rconn, conn)
			closeChan <- 1
		}()
		go func() {
			io.Copy(conn, rconn)
			closeChan <- 1
		}()

		<-closeChan
		rconn.Close()
		conn.Close()
	}
}
