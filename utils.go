package main

import (
	"bufio"
	"strings"
)

// readMsg reads a string from a reader up to newline delimiter
func readMsg(r *bufio.Reader) (string, error) {
	msg, err := r.ReadString('\n')

	if err != nil {
		return "", err
	}
	return strings.TrimSpace(msg), nil
}
