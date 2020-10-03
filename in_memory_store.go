package main

import "fmt"

type InMemoryStore struct {
	records map[string]SIPRecord
}

func (i *InMemoryStore) Find(aor string) (*SIPRecord, error) {
	record, ok := i.records[aor]
	if !ok {
		return nil, fmt.Errorf("No record found for aor %s", aor)
	}
	return &record, nil
}
