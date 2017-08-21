package hooks

import (
	"errors"
	"os/exec"

	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

type ExecCommand struct {
	List []string
	Log  *log.Logger
	Stop bool
}

func NewExecommand(manifest utils.ExecCommandManifest, logger *log.Logger) *ExecCommand {
	return &ExecCommand{
		List: manifest.List,
		Log:  logger,
		Stop: manifest.StopOnError,
	}
}

func (e ExecCommand) PreDeployment(username string, env string, service string, podVersion string) error {
	return e.executeCommands()
}

func (e ExecCommand) PostDeployment(username string, env string, service string, podVersion string) error {
	return e.executeCommands()
}

func (e ExecCommand) executeCommands() error {
	var err error = nil
	for _, commandName := range e.List {
		if _, error := exec.Command(commandName).CombinedOutput(); error != nil {
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
