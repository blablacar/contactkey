package commands

import (
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to a previous version",
}

type Rollback struct {
	Env     string
	Service string
}

func (r Rollback) execute() {
}

func (r Rollback) fill(service string, env string) {
	r.Env = env
	r.Service = service
}
