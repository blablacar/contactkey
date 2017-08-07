package commands

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/remyLemeunier/contactkey/context"
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
	fill(context *context.Context, service string, env string)
	execute()
}

func fill(cck CckCommand, config *utils.Config, service string, env string) {
	filePath := path.Join(config.WorkPath, fmt.Sprintf("%s.yml", service))
	manifestFile, err := utils.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Unable to read file: %q with err: %q", filePath, err)
		os.Exit(1)
	}

	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		fmt.Printf("LoadConfig failed with err %q", err)
		os.Exit(1)
	}

	ctxt, err := context.NewContext(config, manifest)
	if err != nil {
		fmt.Printf("NewContext failed with err %q", err)
		os.Exit(1)
	}

	cck.fill(ctxt, service, env)
}

func execute(cck CckCommand) {
	cck.execute()
}
