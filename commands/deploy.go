package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/spf13/cobra"
)

func init() {
	deployCmd.PersistentFlags().StringVar(&branch, "branch", "", "Change the branch from the default one.")
}

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
	// If the branch is null it will use the default one.
	sha1ToDeploy, err := d.Context.Vcs.RetrieveSha1ForProject(branch)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve version from service(%q): %q", d.Service, err)
		os.Exit(1)
	}

	podVersion, err := d.Context.RepositoryManager.RetrievePodVersion(sha1ToDeploy)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failted to retrieve pod version: %q", err)
		os.Exit(1)
	}

	deployedVersions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve DEPLOYED version from service(%q) in env %q: %q", d.Service, d.Env, err)
		os.Exit(1)
	}

	if podVersion == "" {
		fmt.Fprintf(d.Writer, "We have not found the pod version with the the sha1 %q \n"+
			"The pod has not been created.", sha1ToDeploy)
		os.Exit(1)
	}

	needToDeploy := false
	for _, deployedVersion := range deployedVersions {
		if deployedVersion != sha1ToDeploy {
			needToDeploy = true
		}
	}

	if needToDeploy == false {
		fmt.Fprintf(d.Writer, "I can't deploy(sorry) as the version you want to deploy \n"+
			"is the same as the version deployed (%q)", sha1ToDeploy)
		os.Exit(1)
	}

	fmt.Fprintf(d.Writer, "Going to deploy pod version %q", podVersion)
	for _, hook := range d.Context.Hooks {
		//@TODO Add a logger and log error coming from hooks
		hook.PreDeployment(d.Env, d.Service, podVersion)
	}

	for _, hook := range d.Context.Hooks {
		//@TODO Add a logger and log error coming from hooks
		hook.PostDeployment(d.Env, d.Service, podVersion)
	}
}

func (d *Deploy) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
	d.TableWriter = tablewriter.NewWriter(os.Stdout)
	d.Writer = os.Stdout
}
