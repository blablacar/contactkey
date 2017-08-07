package commands

import (
	"fmt"
	"os"

	"path"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/remyLemeunier/contactkey/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the version of the service",
}

type List struct {
	*context.Context
	Env         string
	Service     string
	Config      *utils.Config
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

func (l *List) fill(config *utils.Config, service string, env string) {
	l.Env = env
	l.Service = service
	l.Config = config
	l.TableWriter = tablewriter.NewWriter(os.Stdout)

	filePath := path.Join(l.Config.WorkPath, fmt.Sprintf("%s.yml", l.Service))
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

	ctxt, err := context.NewContext(l.Config, manifest)
	if err != nil {
		fmt.Printf("NewContext failed with err %q", err)
		os.Exit(1)
	}
	l.Context = ctxt
}
