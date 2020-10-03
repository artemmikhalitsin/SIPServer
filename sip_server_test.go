package main

import (
	"net"
	"testing"
	"time"
)

func TestSIPServer(t *testing.T) {

	t.Run("It should accept TCP connections", func(t *testing.T) {
		port := ":1123"
		server, _ := NewSIPServer(port)
		go server.Listen()

		// Wait briefly for server to start
		time.Sleep(10 * time.Millisecond)

		conn, err := net.Dial("tcp", port)
		if err != nil {
			t.Errorf("Expected to connect to server but got error: %v", err)
		}
		defer conn.Close()
	})
}
