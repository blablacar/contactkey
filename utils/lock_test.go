package utils

import (
	"os"
	"testing"
)

func TestFileLock_Lock(t *testing.T) {
	filePath := os.TempDir() + "/contactkey.lock"
	fileLock := FileLock{
		filePath,
	}

	canLock, err := fileLock.Lock()
	if err != nil {
		t.Fatalf("Error raised %q", err)
	}

	if canLock == false {
		t.Fatal("It should have been unlocked.")
	}

	// Trying to lock a second time it should return false
	canLock, err = fileLock.Lock()
	if canLock == true {
		t.Fatal("It should have been locked by the previous Lock() query.")
	}

	if err := fileLock.Unlock(); err != nil {
		t.Fatal("Impossible to remove file previously created.")
	}
}
