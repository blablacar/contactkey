package hooks

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

type ExecCommand struct {
	OnInit       []utils.CommandList
	OnPreDeploy  []utils.CommandList
	OnPostDeploy []utils.CommandList
	Stop         bool
}

func NewExecommand(manifest utils.ExecCommandManifest) *ExecCommand {
	return &ExecCommand{
		OnInit:       manifest.OnInit,
		OnPreDeploy:  manifest.OnPreDeploy,
		OnPostDeploy: manifest.OnPostDeploy,
		Stop:         manifest.StopOnError,
	}
}

func (e ExecCommand) Init() error {
	return executeCommands(e.OnInit)
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
		log.Debugf("Executing command %s %s", commandName.Command, strings.Join(commandName.Args, " "))
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
