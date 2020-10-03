package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	t.Run("Can retrieve a record by its AOR", func(t *testing.T) {

		store := &InMemoryStore{}
		record1 := SIPRecord{
			AddressOfRecord: "aor1",
		}
		record2 := SIPRecord{
			AddressOfRecord: "aor2",
			TenantId:        "123",
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

	t.Run("Can load records from file", func(t *testing.T) {
		contents := `{"addressOfRecord":"aor1", "tenantId":"123"}
{"addressOfRecord":"aor2", "tenantId":"234"}`
		file, closeFile := createTempFile(t, "records", contents)
		defer closeFile()
		store := InMemoryStoreFromFile(file)

		// Retrieve first record
		got, err := store.Find("aor1")
		want := SIPRecord{
			AddressOfRecord: "aor1",
			TenantId:        "123",
		}
		assertNoError(t, err)
		assertRecordEquals(t, &want, got)

		// Retrieve second record
		want = SIPRecord{
			AddressOfRecord: "aor2",
			TenantId:        "234",
		}
		got, err = store.Find("aor2")
		assertNoError(t, err)
		assertRecordEquals(t, &want, got)
		if len(store.records) != 2 {
			t.Errorf("Expected store to contain 2 records, but it only has %d", len(store.records))
		}
	})
}

func assertRecordEquals(t *testing.T, want, got *SIPRecord) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func createTempFile(t *testing.T, name, data string) (*os.File, func()) {
	t.Helper()
	file, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}
	file.Write([]byte(data))
	file.Seek(0, 0)

	closeFunc := func() {
		file.Close()
	}

	return file, closeFunc
}
