package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	t.Run("Can retrieve a registration by its AOR", func(t *testing.T) {

		store := &InMemoryStore{}
		registration1 := SIPRegistration{
			AddressOfRecord: "aor1",
		}
		registration2 := SIPRegistration{
			AddressOfRecord: "aor2",
			TenantId:        "123",
		}
		var registrations = map[string]SIPRegistration{
			"aor1": registration1,
			"aor2": registration2,
		}
		store.registrations = registrations

		// Try to find registration 1
		want := registration1
		got, err := store.Find("aor1")
		assertNoError(t, err)
		assertRecordEquals(t, &want, got)

		// Try to find registration 2
		want = registration2
		got, err = store.Find("aor2")
		assertNoError(t, err)
		assertRecordEquals(t, &want, got)

		// Try to retrieve non-existing registration
		_, err = store.Find("aor3")
		if err == nil {
			t.Errorf("Exoected error on non-existing registration didn't get one")
		}
	})

	t.Run("Can load registrations from file", func(t *testing.T) {
		contents := `{"addressOfRecord":"aor1", "tenantId":"123"}
{"addressOfRecord":"aor2", "tenantId":"234"}`
		file, closeFile := createTempFile(t, "registrations", contents)
		defer closeFile()
		store := InMemoryStoreFromFile(file)

		// Retrieve first registration
		got, err := store.Find("aor1")
		want := SIPRegistration{
			AddressOfRecord: "aor1",
			TenantId:        "123",
		}
		assertNoError(t, err)
		assertRecordEquals(t, &want, got)

		// Retrieve second registration
		want = SIPRegistration{
			AddressOfRecord: "aor2",
			TenantId:        "234",
		}
		got, err = store.Find("aor2")
		assertNoError(t, err)
		assertRecordEquals(t, &want, got)
		if len(store.registrations) != 2 {
			t.Errorf("Expected store to contain 2 registrations, but it only has %d", len(store.registrations))
		}
	})
}

func assertRecordEquals(t *testing.T, want, got *SIPRegistration) {
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
