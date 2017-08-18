package utils

import (
	"errors"
	"fmt"
	"os"
)

func CheckIfIsLaunchedInAScreen() error {
	term := os.Getenv("TERM")
	if term == "screen" {
		return nil
	}

	return errors.New(fmt.Sprintf("Terminal is :%q instead of screen.", term))
}
