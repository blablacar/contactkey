package main

import (
	"fmt"

	"runtime/debug"

	"github.com/blablacar/contactkey/commands"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{"stacktrace": string(debug.Stack())}).Error(fmt.Sprintf("Recovered from panic %q", r))

		}
	}()
	commands.Execute()
}
