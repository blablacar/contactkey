package commands

import (
	"errors"
	"fmt"
)

var typeRegistry = make(map[string]CckCommand)

func makeInstance(name string) (CckCommand, error) {
	if _, ok := typeRegistry[name]; !ok {
		return nil, errors.New(fmt.Sprintf("Struct not found %s", name))

	}

	return typeRegistry[name], nil
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
