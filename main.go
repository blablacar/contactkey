package main

import (
	"fmt"

	"runtime/debug"

	"github.com/remyLemeunier/contactkey/commands"
	log "github.com/sirupsen/logrus"
	"github.com/wercker/journalhook"
)

func main() {
	log.AddHook(&journalhook.JournalHook{})
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{"stacktrace": string(debug.Stack())}).Error(fmt.Sprintf("Recovered from panic %q", r))

		}
	}()

	commands.Execute()
}
