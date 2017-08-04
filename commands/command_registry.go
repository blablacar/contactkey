package commands

import (
	"errors"
	"fmt"

	"github.com/remyLemeunier/contactkey/utils"
)

var typeRegistry = make(map[string]CckCommand)

func makeInstance(cfg *utils.Config, name string, service string, env string) (CckCommand, error) {
	if _, ok := typeRegistry[name]; !ok {
		return nil, errors.New(fmt.Sprintf("Struct not found %s", name))

	}

	cckCommand := typeRegistry[name]
	fill(cckCommand, cfg, service, env)

	return cckCommand, nil
}

func init() {
	typeRegistry["deploy"] = &Deploy{}
	typeRegistry["diff"] = &Diff{}
	typeRegistry["list"] = &List{}
	typeRegistry["rollback"] = &Rollback{}
}

type CckCommand interface {
	fill(config *utils.Config, service string, env string)
	execute()
}

func fill(cck CckCommand, config *utils.Config, service string, env string) {
	cck.fill(config, service, env)
}

func execute(cck CckCommand) {
	cck.execute()
}
