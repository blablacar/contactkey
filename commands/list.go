package commands

import (
	"github.com/remyLemeunier/contactkey/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the version of the service",
}

type List struct {
	Env     string
	Service string
	Config  *utils.Config
}

func (l List) execute() {

}

func (l List) fill(config *utils.Config, service string, env string) {
	l.Env = env
	l.Service = service
	l.Config = config
}
