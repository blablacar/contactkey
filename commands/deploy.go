package commands

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/spf13/cobra"
)

var userName = "Mister Robot"

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
	currentUser, err := user.Current()
	if err == nil {
		userName = currentUser.Name
	}

	// The lock system is not mandatory
	if d.Context.LockSystem != nil {
		fmt.Fprintf(d.Writer, "Trying to lock the lock command for service %q and env %q \n", d.Service, d.Env)
		canLock, err := d.Context.LockSystem.Lock(d.Env, d.Service)
		if err != nil {
			fmt.Fprintf(d.Writer, "Failed to lock, error raised: %q", err)
		}

		if canLock == false {
			fmt.Fprint(d.Writer, "Another command is currently running")
			return
		}

		defer func(d *Deploy) {
			d.Context.LockSystem.Unlock(d.Env, d.Service)
			if err != nil {
				fmt.Fprintf(d.Writer, "Failed to unlock, error raised: %q", err)
			}
		}(d)
	}

	// If the branch is null it will use the default one.
	sha1ToDeploy, err := d.Context.Vcs.RetrieveSha1ForProject(branch)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve version from service(%q): %q", d.Service, err)
		return
	}

	podVersion, err := d.Context.RepositoryManager.RetrievePodVersion(sha1ToDeploy)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve pod version: %q", err)
		return
	}

	deployedVersions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve DEPLOYED version from service(%q) in env %q: %q", d.Service, d.Env, err)
		return
	}

	if podVersion == "" {
		fmt.Fprintf(d.Writer, "We have not found the pod version with the the sha1 %q \n"+
			"The pod has not been created.", sha1ToDeploy)
		return
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
		return
	}

	fmt.Fprintf(d.Writer, "Going to deploy pod version %q", podVersion)
	for _, hook := range d.Context.Hooks {
		//@TODO Add a logger and log error coming from hooks
		hook.PreDeployment(userName, d.Env, d.Service, podVersion)
	}

	if err := d.Context.Deployer.Deploy(d.Env, podVersion); err != nil {
		fmt.Fprintf(d.Writer, "Failed to deploy : %q", err)
	}

	for _, hook := range d.Context.Hooks {
		//@TODO Add a logger and log error coming from hooks
		hook.PostDeployment(userName, d.Env, d.Service, podVersion)
	}
}

func (d *Deploy) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
	d.TableWriter = tablewriter.NewWriter(os.Stdout)
	d.Writer = os.Stdout
}
