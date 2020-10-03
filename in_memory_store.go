package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// InMemoryStore stores a table of SIP registrations in memory
type InMemoryStore struct {
	registrations map[string]SIPRegistration
}

// NewInMemoryStore creates a new empty InMemoryStore
func NewInMemoryStore() *InMemoryStore {
	registrations := make(map[string]SIPRegistration)
	return &InMemoryStore{
		registrations: registrations,
	}
}

// InMemoryStoreFromFile creates a InMemoryStore
// and populates it with registrations from the supplied file
func InMemoryStoreFromFile(file *os.File) *InMemoryStore {
	store := NewInMemoryStore()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var registration SIPRegistration
		line := scanner.Text()
		json.Unmarshal([]byte(line), &registration)

		store.registrations[registration.AddressOfRecord] = registration
	}

	return store
}

// Find retrieves a registration given an address of record
func (i *InMemoryStore) Find(aor string) (*SIPRegistration, error) {
	registration, ok := i.registrations[aor]
	if !ok {
		return nil, fmt.Errorf("No registration found for address of record %s", aor)
	}
	return &registration, nil
}
