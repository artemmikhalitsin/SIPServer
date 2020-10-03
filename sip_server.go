package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const connectionClosedMessage = "Connection closed."

// An SIPStore is a store that holds SIP records
type SIPStore interface {
	Find(aor string) (*SIPRecord, error)
}

// SIPRecordServer accepts incoming TCP connections and processess requests
// for SIP record lookups
type SIPRecordServer struct {
	address  string
	listener net.Listener
	store    SIPStore
	timeout  time.Duration
	close    chan interface{} // signals whether the server has stopped accepting connections
	wg       sync.WaitGroup
}

// NewSIPRecordServer creates a new SIP Server with the given parameters
func NewSIPRecordServer(address string, store SIPStore, timeout time.Duration) (*SIPRecordServer, chan interface{}) {
	closed := make(chan interface{})
	return &SIPRecordServer{
		address: address,
		store:   store,
		timeout: timeout,
		close:   closed,
	}, closed
}

// Listen starts the server's TCP listener and begins accepting incoming connections
func (s *SIPRecordServer) Listen() error {
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
				log.Printf("Error accepting connection: %v", err)
			}
		} else {
			go s.serveConnection(conn)
		}
	}
}

// Close tells the server to stop accepting connections, closes its listener and
// then waits for any ongoing connections to terminate
func (s *SIPRecordServer) Close() {
	close(s.close)
	s.listener.Close()
	s.wg.Wait()
}

// serveConnection registers and handles incoming connections as long as they
// stay active. Once a connection is inactive for longer than specified
// in the server configuration, it is terminated
func (s *SIPRecordServer) serveConnection(conn net.Conn) {
	s.wg.Add(1)
	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(s.timeout))
		msg, err := readMsg(reader)
		if err != nil {
			fmt.Fprintln(conn, connectionClosedMessage)
			conn.Close()
			s.wg.Done()
			return
		}
		record, err := s.store.Find(msg)
		if err != nil {
			fmt.Fprintf(conn, "Could not find record with address %s\n", msg)
		} else {
			response, _ := json.Marshal(record)
			fmt.Fprintln(conn, string(response))
		}
	}
}
