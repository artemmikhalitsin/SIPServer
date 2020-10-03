package main

import (
	"bufio"
	"fmt"
	"net"
)

type SIPServer struct {
	address  string
	listener net.Listener
	store    SIPStore
}

func NewSIPServer(port string, store SIPStore) *SIPServer {
	return &SIPServer{
		address: port,
		store:   store,
	}
}

func (s *SIPServer) Listen() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("Error listening on %s: %w", s.address, err)
	}
	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go s.lookupRecords(conn)
	}
}

func (s *SIPServer) lookupRecords(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		s.store.Find(msg)
		conn.Write([]byte("Found\n"))
	}
}

type SIPStore interface {
	Find(aor string)
}
