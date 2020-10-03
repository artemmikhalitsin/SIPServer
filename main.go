package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	regsFile, err := os.Open("regs")
	if err != nil {
		log.Fatalf("Error opening the SIP registration file: %v", err)
	}
	store := InMemoryStoreFromFile(regsFile)
	regsFile.Close()
	timeout := time.Second * 10
	address := ":4444"
	server, closed := NewSIPRecordServer(address, store, timeout)
	fmt.Println("Accepting TCP connection at address", address)
	go server.Listen()
	<-closed
	log.Println("Server terminated")
}

type echoStore struct{}

func (e *echoStore) Find(aor string) string {
	return aor
}
