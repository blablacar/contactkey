package commands

import (
	"fmt"

	"github.com/remyLemeunier/contactkey/context"
	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

var typeRegistry = make(map[string]CckCommand)

func makeInstance(cfg *utils.Config, name string, service string, env string, filePath string) (CckCommand, error) {
	if _, ok := typeRegistry[name]; !ok {
		return nil, fmt.Errorf("Struct not found %s", name)

	}

	cckCommand := typeRegistry[name]
	err := fill(cckCommand, cfg, service, env, filePath)
	if err != nil {
		return nil, err
	}

	return cckCommand, nil
}

func init() {
	typeRegistry["deploy"] = &Deploy{}
	typeRegistry["diff"] = &Diff{}
	typeRegistry["list"] = &List{}
	typeRegistry["rollback"] = &Rollback{}
}

type CckCommand interface {
	fill(context *context.Context, service string, env string)
	execute()
}

func fill(cck CckCommand, config *utils.Config, service string, env string, filePath string) error {
	manifestFile, err := utils.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Unable to read file: %q with err: %q", filePath, err)
	}

	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		return fmt.Errorf("LoadConfig failed with err %q", err)
	}

	ctxt, err := context.NewContext(config, manifest)
	if err != nil {
		return fmt.Errorf("NewContext failed with err %q", err)
	}

	for _, hook := range ctxt.Hooks {
		err = hook.Init()
		if hook.StopOnError() == true && err != nil {
			return fmt.Errorf("Init hook failed: %q", err)
		} else if err != nil {
			log.Debugf("Init hook failed: %q", err)
		}
	}

	cck.fill(ctxt, service, env)

	return nil
}

func execute(cck CckCommand) {
	cck.execute()
}
