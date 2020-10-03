package main

import (
	"bufio"
	"strings"
)

// readMessage reads messages from a reader up to newline
// and returns it over a channel
func readMessage(r *bufio.Reader) <-chan string {
	msgChannel := make(chan string, 1)
	go func() {
		msg, _ := r.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if msg != "" {
			msgChannel <- msg
		}
	}()
	return msgChannel
}

func readMsg(r *bufio.Reader) (string, error) {
	msg, err := r.ReadString('\n')

	if err != nil {
		return "", err
	}
	return strings.TrimSpace(msg), nil
}
