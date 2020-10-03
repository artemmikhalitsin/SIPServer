package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const connectionClosedMessage = "Connection closed."

type SIPStore interface {
	Find(aor string)
}

type SIPServer struct {
	address  string
	listener net.Listener
	store    SIPStore
	config   *ServerConfig
	close    chan interface{}
	wg       sync.WaitGroup
}

type ServerConfig struct {
	timeout time.Duration
}

func NewSIPServer(port string, store SIPStore, config *ServerConfig) *SIPServer {
	return &SIPServer{
		address: port,
		store:   store,
		config:  config,
		close:   make(chan interface{}),
	}
}

func (s *SIPServer) Listen() error {
	listener, err := net.Listen("tcp", s.address)
	s.listener = listener
	if err != nil {
		return fmt.Errorf("Error listening on %s: %w", s.address, err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-s.close:
				return nil
			default:
				log.Printf("Error accepting connection: %v\n", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.serveConnection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *SIPServer) Close() {
	close(s.close)
	s.listener.Close()
	s.wg.Wait()
}

func (s *SIPServer) serveConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		select {
		case msg := <-readMessage(reader):
			s.store.Find(msg)
			conn.Write([]byte("Found\n"))
		case <-time.After(s.config.timeout):
			conn.Write([]byte(connectionClosedMessage + "\n"))
			conn.Close()
			return
		}
	}
}
