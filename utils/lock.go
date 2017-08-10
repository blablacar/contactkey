package utils

import (
	"io/ioutil"
	"os"
)

type Lock interface {
	Lock() (bool, error)
	Unlock() error
}

type FileLock struct {
	FilePath string
}

// The lock function will also verify if
// it's already locked or not.
// Return false if it can't lock true otherwise
func (f FileLock) Lock() (bool, error) {
	// Check if the file exists
	if _, err := os.Stat(f.FilePath); !os.IsNotExist(err) {
		return false, err
	}
	// 0222 = --w--w--w-
	err := ioutil.WriteFile(f.FilePath, []byte("locked by contactKey"), 0222)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Check if the file exists and remove it
// if necessary.
func (f FileLock) Unlock() error {
	if _, err := os.Stat(f.FilePath); !os.IsExist(err) {
		return nil
	}

	if err := os.Remove(f.FilePath); err != nil {
		return err
	}

	return nil
}
