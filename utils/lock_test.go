package utils

import (
	"os"
	"testing"
)

func TestFileLock_Lock(t *testing.T) {
	fileLock := FileLock{
		os.TempDir(),
	}
	preprodEnv := "preprod"
	prodEnv := "prod"
	service := "webooks"
	canLock, err := fileLock.Lock(preprodEnv, service)
	if err != nil {
		t.Fatalf("Error raised %q", err)
	}

	if canLock == false {
		t.Fatal("It should have been unlocked.")
	}

	// Trying we another env it should be able to lock
	canLock, err = fileLock.Lock(prodEnv, service)
	if err != nil {
		t.Fatalf("Error raised %q", err)
	}

	if canLock == false {
		t.Fatal("It should have been unlocked.")
	}

	// Trying to lock a second time it should return false
	canLock, err = fileLock.Lock(preprodEnv, service)
	if canLock == true {
		t.Fatal("It should have been locked by the previous Lock() query.")
	}

	if err := fileLock.Unlock(preprodEnv, service); err != nil {
		t.Fatal("Impossible to remove file previously created.")
	}

	if err := fileLock.Unlock(prodEnv, service); err != nil {
		t.Fatal("Impossible to remove file previously created.")
	}
}
