package commands

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the version of the service",
}

type List struct {
	Context     *context.Context
	Env         string
	Service     string
	TableWriter *tablewriter.Table
}

func (l List) execute() {
	versions, err := l.Context.Deployer.ListVersions(l.Env)
	if err != nil {
		fmt.Printf("Failed to list versions : %s", err)
		os.Exit(1)
	}

	l.TableWriter.SetHeader([]string{"instance", "version"})
	for i, v := range versions {
		l.TableWriter.Append([]string{i, v})
	}
	l.TableWriter.Render()

}

func (l *List) fill(context *context.Context, service string, env string) {
	l.Env = env
	l.Service = service
	l.Context = context
	l.TableWriter = tablewriter.NewWriter(os.Stdout)
}
