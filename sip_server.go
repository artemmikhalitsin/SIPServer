package main

import (
	"fmt"
	"net"
)

type SIPServer struct {
	address  string
	listener net.Listener
}

func NewSIPServer(port string) (*SIPServer, error) {
	return &SIPServer{
		address: port,
	}, nil
}

func (s *SIPServer) Listen() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("Error listening on %s: %w", s.address, err)
	}
	defer listener.Close()
	for {
		listener.Accept()
	}
}
