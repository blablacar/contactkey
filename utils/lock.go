package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var filePathPattern = "%s/contactkey_%s_%s.lock"

type Lock interface {
	Lock(env string, service string) (bool, error)
	Unlock(env string, service string) error
}

func NewFileLock(config FileLockConfig) (*FileLock, error) {
	if config.FilePath == "" {
		return nil, errors.New("You need to define a filePath for fileLock in the config.")
	}

	return &FileLock{
		config.FilePath,
	}, nil
}

type FileLock struct {
	FileDir string
}

// The lock function will also verify if
// it's already locked or not.
// Return false if it can't lock true otherwise
func (f FileLock) Lock(env string, service string) (bool, error) {
	// Check if the file exists
	filePath := fmt.Sprintf(filePathPattern, f.FileDir, env, service)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return false, err
	}
	// 0222 = --w--w----
	err := ioutil.WriteFile(filePath, []byte("locked by contactKey"), 0222)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Check if the file exists and remove it
// if necessary.
func (f FileLock) Unlock(env string, service string) error {
	filePath := fmt.Sprintf(filePathPattern, f.FileDir, env, service)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
