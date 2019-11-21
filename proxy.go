package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Proxy needs to be implemented by any proxy
type Proxy interface {
	StartProxy() error
}

// TCPProxy represents a TCP Proxy
type TCPProxy struct {
	Src *net.TCPAddr
	Dst *net.TCPAddr
}

// NewProxy returns a new proxy instance according to proxy type
func NewProxy(proxyType, src, dst string) (Proxy, error) {
	switch proxyType {
	case "tcp":
		return NewTCPProxy(src, dst)
	default:
		return nil, fmt.Errorf("unrecognized proxy type: %s", proxyType)
	}
}

// NewTCPProxy returns a new proxy instance
func NewTCPProxy(src, dst string) (*TCPProxy, error) {
	laddr, err := net.ResolveTCPAddr("tcp", cfg.Source)
	if err != nil {
		return nil, err
	}

	raddr, err := net.ResolveTCPAddr("tcp", cfg.Destination)
	if err != nil {
		return nil, err
	}

	return &TCPProxy{
		Src: laddr,
		Dst: raddr,
	}, nil
}

// StartProxy starts a TCP Proxy
func (p *TCPProxy) StartProxy() error {
	listener, err := net.ListenTCP("tcp", p.Src)
	if err != nil {
		return fmt.Errorf("error while listening: %v", err)
	}

	for {
		conn, err := listener.AcceptTCP()
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
