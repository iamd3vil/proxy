package main

import (
	"net"
	"sync"
)

// BUFFERSIZE is the maximum length of the packet received
const BUFFERSIZE = 2000

type proxyConn struct {
	clientAddr *net.UDPAddr
	destConn   *net.UDPConn
	serverConn *net.UDPConn
}

// UDPProxy implements Proxy interface for proxying UDP
type UDPProxy struct {
	sync.Mutex
	ServerAddr *net.UDPAddr
	ServerConn *net.UDPConn
	DestAddr   *net.UDPAddr
	clients    map[string]proxyConn
}

// NewUDPProxy returns a new UDP Proxy
func NewUDPProxy(serverAddr, destAddr string) (Proxy, error) {
	sAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	dAddr, err := net.ResolveUDPAddr("udp", destAddr)
	if err != nil {
		return nil, err
	}

	return &UDPProxy{
		ServerAddr: sAddr,
		DestAddr:   dAddr,
		clients:    make(map[string]proxyConn),
	}, nil
}

// StartProxy starts a UDP Proxy
func (p *UDPProxy) StartProxy() error {
	// Listen to UDP packets
	sConn, err := net.ListenUDP("udp", p.ServerAddr)
	if err != nil {
		return err
	}

	p.ServerConn = sConn

	return p.readFromServerConn()
}

// Starts to read from the UDP proxy server connection
func (p *UDPProxy) readFromServerConn() error {
	packet := make([]byte, BUFFERSIZE)
	for {
		n, clientAddr, err := p.ServerConn.ReadFromUDP(packet)
		if err != nil {
			return err
		}

		// Check if we have clientAddr in the clients map
		p.Lock()
		pConn, ok := p.clients[clientAddr.String()]
		if !ok {
			// Start a proxy connection for this client
			prConn := proxyConn{}
			prConn.clientAddr = clientAddr
			prConn.serverConn = p.ServerConn

			// Make a connection to dest
			dConn, err := net.DialUDP("udp", nil, p.DestAddr)
			if err != nil {
				return err
			}
			prConn.destConn = dConn
			p.clients[clientAddr.String()] = prConn
			p.Unlock()

			// Start a goroutine to listen from the dest conn and write packets to the client
			go prConn.copyFromDestConnToClient()

			// Write to the server connection
			_, err = prConn.destConn.Write(packet[:n])
			if err != nil {
				return err
			}
			continue
		}
		p.Unlock()
		_, err = pConn.destConn.Write(packet[:n])
		if err != nil {
			return err
		}
	}
}

func (pc *proxyConn) copyFromDestConnToClient() error {
	packet := make([]byte, BUFFERSIZE)
	for {
		n, err := pc.destConn.Read(packet)
		if err != nil {
			return err
		}

		_, err = pc.serverConn.WriteToUDP(packet[:n], pc.clientAddr)
		if err != nil {
			return err
		}
	}
}
