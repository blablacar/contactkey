package commands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the service in an environment",
}

type Deploy struct {
	Env         string
	Service     string
	Context     *context.Context
	TableWriter *tablewriter.Table
	Writer      io.Writer
}

func (d *Deploy) execute() {
	serviceVersionFromPod, err := d.Context.RepositoryManager.RetrieveServiceVersionFromPod()
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve service version from Pod: %q", err)
		os.Exit(1)
	}

	sha1ToDeploy, err := d.Context.Vcs.RetrieveSha1ForProject("")
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve version from service(%q): %q", d.Service, err)
		os.Exit(1)
	}

	deployedVersions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve DEPLOYED version from service(%q) in env %q: %q", d.Service, d.Env, err)
		os.Exit(1)
	}

	// It means that the pod has not been created
	if !strings.Contains(sha1ToDeploy, serviceVersionFromPod) {
		fmt.Fprintf(d.Writer, "The version to deploy(%q) differs from the pod version (%q). \n"+
			"The pod has not been created.", sha1ToDeploy, serviceVersionFromPod)
		os.Exit(1)
	}

	needToDeploy := false
	for _, deployedVersion := range deployedVersions {
		if deployedVersion != sha1ToDeploy {
			needToDeploy = true
		}
	}

	if needToDeploy == false {
		fmt.Fprintf(d.Writer, "The version to deploy(%q) differs from the pod version (%q). \n"+
			"The pod has not been created.", sha1ToDeploy, serviceVersionFromPod)
		os.Exit(1)
	}

	fmt.Fprintf(d.Writer, "Going to deploy version %q", sha1ToDeploy)
}

func (d *Deploy) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
	d.TableWriter = tablewriter.NewWriter(os.Stdout)
	d.Writer = os.Stdout
}
