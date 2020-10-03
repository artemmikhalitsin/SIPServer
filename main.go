package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	address := fmt.Sprintf(":%s", os.Getenv("SIP_PORT"))
	if address == ":" {
		// default port
		address = ":4444"
	}
	regsFilePath := os.Getenv("SIP_REGS_FILE")
	if regsFilePath == "" {
		regsFilePath = "regs"
	}

	regsFile, err := os.Open(regsFilePath)
	if err != nil {
		log.Fatalf("Error opening the SIP registration file: %v", err)
	}
	store := InMemoryStoreFromFile(regsFile)
	regsFile.Close()
	timeout := time.Second * 10
	server, closed := NewSIPRecordServer(address, store, timeout)
	log.Println("Accepting TCP connection at address", address)
	go server.Listen()
	<-closed
	log.Println("Server terminated")
}
