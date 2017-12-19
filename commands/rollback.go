package commands

import (
	"github.com/remyLemeunier/contactkey/context"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to a previous version",
}

type Rollback struct {
	Env     string
	Service string
	Context *context.Context
}

func (r *Rollback) execute() error {
	return nil
}

func (r *Rollback) fill(context *context.Context, service string, env string) {
	r.Env = env
	r.Service = service
	r.Context = context
}
