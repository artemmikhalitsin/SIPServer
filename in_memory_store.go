package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type InMemoryStore struct {
	registrations map[string]SIPRegistration
}

func NewInMemoryStore() *InMemoryStore {
	registrations := make(map[string]SIPRegistration)
	return &InMemoryStore{
		registrations: registrations,
	}
}

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

func (i *InMemoryStore) Find(aor string) (*SIPRegistration, error) {
	registration, ok := i.registrations[aor]
	if !ok {
		return nil, fmt.Errorf("No registration found for address of record %s", aor)
	}
	return &registration, nil
}
