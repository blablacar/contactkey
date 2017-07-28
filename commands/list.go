package commands

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the version of the service",
}

type List struct {
	Env     string
	Service string
}

func (l List) execute() {
}

func (l List) fill(service string, env string) {
	l.Env = env
	l.Service = service
}
