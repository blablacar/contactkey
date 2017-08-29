package hooks

import (
	"errors"
	"os/exec"

	"github.com/remyLemeunier/contactkey/utils"
)

type ExecCommand struct {
	OnPreDeploy  []utils.CommandList
	OnPostDeploy []utils.CommandList
	Stop         bool
}

func NewExecommand(manifest utils.ExecCommandManifest) *ExecCommand {
	return &ExecCommand{
		OnPreDeploy:  manifest.OnPreDeploy,
		OnPostDeploy: manifest.OnPostDeploy,
		Stop:         manifest.StopOnError,
	}
}

func (e ExecCommand) PreDeployment(username string, env string, service string, podVersion string) error {
	return executeCommands(e.OnPreDeploy)
}

func (e ExecCommand) PostDeployment(username string, env string, service string, podVersion string) error {
	return executeCommands(e.OnPostDeploy)
}

func executeCommands(commandList []utils.CommandList) error {
	var err error = nil
	for _, commandName := range commandList {
		if _, error := exec.Command(commandName.Command, commandName.Args...).CombinedOutput(); error != nil {
			if err != nil {
				err = errors.New(err.Error() + error.Error())
			} else {
				err = error
			}
		}
	}

	return err
}

func (e ExecCommand) StopOnError() bool {
	return e.Stop
}
