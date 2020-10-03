package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestSIPServer(t *testing.T) {

	t.Run("It should accept TCP connections and process messages", func(t *testing.T) {
		port := ":1123"
		store := &SpyStore{}
		server := NewSIPRecordServer(port, store, time.Millisecond*20)
		defer server.Close()
		go server.Listen()

		// Wait briefly for server to start
		time.Sleep(10 * time.Millisecond)

		conn, err := net.Dial("tcp", port)
		if err != nil {
			t.Errorf("Expected to connect to server but got error: %v", err)
		}
		defer conn.Close()

		responseReader := bufio.NewReader(conn)

		aor := "0142e2fa3543cb32bf000100620002"
		aorMessage := aor + "\n"
		conn.Write([]byte(aorMessage))

		// Wait briefly server to receive message
		time.Sleep(10 * time.Millisecond)

		if store.lastAor != aor {
			t.Errorf("Expected server to look up %q, but got %q instead", aor, store.lastAor)
		}

		select {
		case response := <-readMessage(responseReader):
			if len(response) == 0 {
				t.Errorf("Expected a response from the server, but didn't get one")
			}
			if response == connectionClosedMessage {
				t.Errorf("Got a connection closed message, but didn't expect one")
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Timed out after waiting for a response")
		}
	})

	t.Run("It should close inactive connections", func(t *testing.T) {
		port := ":1123"
		store := &FakeStore{}
		server := NewSIPRecordServer(port, store, time.Millisecond*20)
		defer server.Close()
		go server.Listen()

		// Wait briefly for server to start
		time.Sleep(10 * time.Millisecond)

		conn, err := net.Dial("tcp", port)
		if err != nil {
			t.Errorf("Expected to connect to server but got error: %v", err)
		}
		defer conn.Close()
		responseReader := bufio.NewReader(conn)

		time.Sleep(30 * time.Millisecond)

		select {
		case msg := <-readMessage(responseReader):
			if msg != connectionClosedMessage {
				t.Errorf("Expected a connection closed message, but got %q", msg)
			}
		case <-time.After(50 * time.Millisecond):
			t.Errorf("Client timed out waiting for a response")
		}

	})
}

type FakeStore struct{}

func (f *FakeStore) Find(aor string) {}

type SpyStore struct {
	lastAor string
}

func (s *SpyStore) Find(aor string) {
	s.lastAor = aor
}
