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
		server, _ := NewSIPRecordServer(port, store, time.Millisecond*20)
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

		assertStoreLookups(t, store, 1)
		assertLastAorLookedUp(t, store, aor)

		response, err := readWithTimeout(t, responseReader)
		assertValidResponse(t, response, err)

		// Should serve additional lookups on same connection
		nextAor := "014fb44ecd123b6706000100620005"
		fmt.Fprintln(conn, nextAor)

		time.Sleep(10 * time.Millisecond)

		assertStoreLookups(t, store, 2)
		assertLastAorLookedUp(t, store, nextAor)

		response, err = readWithTimeout(t, responseReader)
		assertValidResponse(t, response, err)
	})

	t.Run("It should process lookups from multiple connections", func(t *testing.T) {
		port := ":1123"
		store := &SpyStore{}
		server, _ := NewSIPRecordServer(port, store, time.Millisecond*20)
		defer server.Close()
		go server.Listen()

		// Wait briefly for server to start
		time.Sleep(10 * time.Millisecond)

		conn1, _, closeConn1 := connectToServer(t, server.address)
		defer closeConn1()
		conn2, _, closeConn2 := connectToServer(t, server.address)
		defer closeConn2()

		sip1 := SIPRecord{
			AddressOfRecord: "aor1",
		}
		sip2 := SIPRecord{
			AddressOfRecord: "aor2",
		}

		fmt.Fprintln(conn1, sip1.AddressOfRecord)
		fmt.Fprintln(conn2, sip2.AddressOfRecord)

		// Wait for server to receive messages
		time.Sleep(10 * time.Millisecond)

		assertStoreLookups(t, store, 2)
		assertAorLookedUp(t, store, "aor1")
		assertAorLookedUp(t, store, "aor2")
	})

	t.Run("It should close inactive connections", func(t *testing.T) {
		port := ":1123"
		store := &FakeStore{}
		server, _ := NewSIPRecordServer(port, store, time.Millisecond*20)
		defer server.Close()
		go server.Listen()

		// Wait briefly for server to start
		time.Sleep(10 * time.Millisecond)

		_, responseReader, closeConn := connectToServer(t, server.address)
		defer closeConn()

		time.Sleep(30 * time.Millisecond)

		msg, err := readWithTimeout(t, responseReader)
		if msg != connectionClosedMessage {
			t.Errorf("Expected a connection closed message, but got %q", msg)
		}
		if err != nil {
			t.Errorf("Error reading response: %w", err)
		}
	})
}

func assertStoreLookups(t *testing.T, store *SpyStore, wanted int) {
	t.Helper()
	if store.lookups != wanted {
		t.Errorf("Expected server to do %d lookups, but got %d", wanted, store.lookups)
	}
}

func assertValidResponse(t *testing.T, response string, err error) {
	t.Helper()
	if len(response) == 0 {
		t.Errorf("Expected a response from the server, but didn't get one")
	}
	if err != nil {
		t.Errorf("Error reading response: %w", err)
	}
	if response == connectionClosedMessage {
		t.Errorf("Got a connection closed message, but didn't expect one")
	}
}

func assertLastAorLookedUp(t *testing.T, store *SpyStore, wanted string) {
	t.Helper()
	if store.lastAor != wanted {
		t.Errorf("Expected server to look up %q, but got %q instead", wanted, store.lastAor)
	}
}

func assertAorLookedUp(t *testing.T, store *SpyStore, aor string) {
	for _, v := range store.aorsLookedUp {
		if v == aor {
			return
		}
	}
	t.Errorf("Expected %s to be looked up but it wasn't", aor)
}

func readWithTimeout(t *testing.T, reader *bufio.Reader) (string, error) {
	msgChannel := make(chan string, 1)
	errorChannel := make(chan error, 1)
	go func() {
		msg, err := readMsg(reader)
		if err != nil {
			errorChannel <- err
		} else {
			msgChannel <- msg
		}
	}()

	select {
	case msg := <-msgChannel:
		return msg, nil
	case err := <-errorChannel:
		return "", err
	case <-time.After(100 * time.Millisecond):
		return "", fmt.Errorf("Timed out waiting for a response")
	}
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

func (f *FakeStore) Find(aor string) (*SIPRecord, error) {
	return nil, nil
}

type SpyStore struct {
	lastAor      string
	aorsLookedUp []string
	lookups      int
}

func (s *SpyStore) Find(aor string) (*SIPRecord, error) {
	s.lastAor = aor
	s.aorsLookedUp = append(s.aorsLookedUp, aor)
	s.lookups++
	return &SIPRecord{}, nil
}
