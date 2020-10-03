package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type InMemoryStore struct {
	records map[string]SIPRecord
}

func NewInMemoryStore() *InMemoryStore {
	records := make(map[string]SIPRecord)
	return &InMemoryStore{
		records: records,
	}
}

func InMemoryStoreFromFile(file *os.File) *InMemoryStore {
	store := NewInMemoryStore()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record SIPRecord
		line := scanner.Text()
		json.Unmarshal([]byte(line), &record)

		store.records[record.AddressOfRecord] = record
	}

	return store
}

func (i *InMemoryStore) Find(aor string) (*SIPRecord, error) {
	record, ok := i.records[aor]
	if !ok {
		return nil, fmt.Errorf("No record found for aor %s", aor)
	}
	return &record, nil
}
