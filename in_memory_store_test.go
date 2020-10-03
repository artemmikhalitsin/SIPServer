package main

import (
	"reflect"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	t.Run("Can retrieve a record by its AOR", func(t *testing.T) {

		store := &InMemoryStore{}
		record1 := SIPRecord{
			addressOfRecord: "aor1",
		}
		record2 := SIPRecord{
			addressOfRecord: "aor2",
			tenantId:        "123",
		}
		var records = map[string]SIPRecord{
			"aor1": record1,
			"aor2": record2,
		}
		store.records = records

		// Try to find record 1
		want := record1
		got, err := store.Find("aor1")
		assertNoError(t, err)
		assertRecordEquals(t, &want, got)

		// Try to find record 2
		want = record2
		got, err = store.Find("aor2")
		assertNoError(t, err)
		assertRecordEquals(t, &want, got)

		// Try to retrieve non-existing record
		_, err = store.Find("aor3")
		if err == nil {
			t.Errorf("Exoected error on non-existing record didn't get one")
		}
	})
}

func assertRecordEquals(t *testing.T, want, got *SIPRecord) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
