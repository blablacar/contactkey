package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the version of the service",
}

type List struct {
	Env         string
	Service     string
	Config      *utils.Config
	TableWriter *tablewriter.Table
}

func (l List) execute() {
	deployFile := path.Join(l.Config.WorkPath, fmt.Sprintf("%s.yml", l.Service))

	manifest, err := utils.LoadDeployfile(&l.Config.DeployDefaults, deployFile)
	if err != nil {
		fmt.Printf("Unexpected error : %q", err)
		os.Exit(1)
	}

	deployer, err := deployers.MakeInstance(manifest.Deploy.Method)
	if err != nil {
		fmt.Printf("Unexpected deployment method : %q", err)
		os.Exit(1)
	}

	versions, err := deployer.ListVersions(l.Env)
	if err != nil {
		fmt.Printf("Failed to list versions with error %q", err)
		os.Exit(1)
	}

	l.TableWriter.SetHeader([]string{"instance", "version"})
	for i, v := range versions {
		l.TableWriter.Append([]string{i, v})
	}
	l.TableWriter.Render()

}

func (l *List) fill(config *utils.Config, service string, env string) {
	l.Env = env
	l.Service = service
	l.Config = config
	l.TableWriter = tablewriter.NewWriter(os.Stdout)
}
