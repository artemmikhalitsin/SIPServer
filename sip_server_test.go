package main

import (
	"bufio"
	"fmt"
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

		conn, responseReader, closeConn := connectToServer(t, server.address)
		defer closeConn()

		aor := "0142e2fa3543cb32bf000100620002"
		fmt.Fprintln(conn, aor)

		// Wait briefly server to receive message
		time.Sleep(10 * time.Millisecond)

		if store.lookups != 1 {
			t.Errorf("Expected server to do %d lookups, but got %d", 1, store.lookups)
		}
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

		// Should serve additional lookups on same connection
		nextAor := "014fb44ecd123b6706000100620005"
		fmt.Fprintln(conn, nextAor)

		time.Sleep(10 * time.Millisecond)

		if store.lookups != 2 {
			t.Errorf("Expected server to do %d lookups, but got %d", 2, store.lookups)
		}
		if store.lastAor != nextAor {
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

		_, responseReader, closeConn := connectToServer(t, server.address)
		defer closeConn()

		time.Sleep(30 * time.Millisecond)

		select {
		case msg := <-readMessage(responseReader):
			if msg != connectionClosedMessage {
				t.Errorf("Expected a connection closed message, but got %q", msg)
			}
		case <-time.After(50 * time.Millisecond):
			t.Errorf("Timed out after waiting for a response")
		}

	})
}

func connectToServer(t *testing.T, address string) (net.Conn, *bufio.Reader, func()) {
	t.Helper()
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Expected to connect to server but got error: %v", err)
	}
	closeFunc := func() {
		conn.Close()
	}
	responseReader := bufio.NewReader(conn)

	return conn, responseReader, closeFunc
}

type FakeStore struct{}

func (f *FakeStore) Find(aor string) {}

type SpyStore struct {
	lastAor string
	lookups int
}

func (s *SpyStore) Find(aor string) {
	s.lastAor = aor
	s.lookups++
}
