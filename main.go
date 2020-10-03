package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	store := &echoStore{}
	timeout := time.Second * 10
	address := ":4444"
	server, closed := NewSIPRecordServer(address, store, timeout)
	fmt.Println("Starting to listen at address", address)
	go server.Listen()
	<-closed
	log.Println("Server terminated")
}

type echoStore struct{}

func (e *echoStore) Find(aor string) string {
	return aor
}
