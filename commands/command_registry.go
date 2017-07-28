package commands

import (
	"github.com/labstack/gommon/log"
)

var typeRegistry = make(map[string]CckCommand)

func makeInstance(name string) CckCommand {
	if _, ok := typeRegistry[name]; !ok {
		log.Fatalf("Struct %s not found", name)
	}

	return typeRegistry[name]
}

func init() {
	typeRegistry["deploy"] = Deploy{}
	typeRegistry["diff"] = Diff{}
	typeRegistry["list"] = List{}
	typeRegistry["rollback"] = Rollback{}
}

type CckCommand interface {
	fill(service string, env string)
	execute()
}

func fill(cck CckCommand, service string, env string) {
	cck.fill(service, env)
}

func execute(cck CckCommand) {
	cck.execute()
}
