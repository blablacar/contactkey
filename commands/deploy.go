package commands

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/utils"
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

	if err := utils.CheckIfIsLaunchedInAScreen(); err != nil && d.Context.ScreenMandatory == true {
		fmt.Fprintf(d.Writer, "Screen error raised: %q \n", err)
		return
	}

	// The lock system is not mandatory
	if d.Context.LockSystem != nil {
		fmt.Fprintf(d.Writer, "Trying to lock the lock command for service %q and env %q \n", d.Service, d.Env)
		canLock, err := d.Context.LockSystem.Lock(d.Env, d.Service)
		if err != nil {
			fmt.Fprintf(d.Writer, "Failed to lock, error raised: %q \n", err)
		}

		if canLock == false {
			fmt.Fprintln(d.Writer, "Another command is currently running")
			return
		}

		defer func(d *Deploy) {
			d.Context.LockSystem.Unlock(d.Env, d.Service)
			if err != nil {
				fmt.Fprintf(d.Writer, "Failed to unlock, error raised: %q \n", err)
			}
		}(d)
	}

	// If the branch is null it will use the default one.
	sha1ToDeploy, err := d.Context.Vcs.RetrieveSha1ForProject(branch)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve version from service(%q): %q \n", d.Service, err)
		return
	}

	podVersion, err := d.Context.Binaries.RetrievePodVersion(sha1ToDeploy)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve pod version: %q \n", err)
		return
	}

	//deployedVersions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	//if err != nil {
	//	fmt.Fprintf(d.Writer, "Failed to retrieve DEPLOYED version from service(%q) in env %q: %q", d.Service, d.Env, err)
	//	return
	//}

	if podVersion == "" {
		fmt.Fprintf(d.Writer, "We have not found the pod version with the the sha1 %q \n"+
			"The pod has not been created. \n", sha1ToDeploy)
		return
	}

	//needToDeploy := false
	//for _, deployedVersion := range deployedVersions {
	//	if deployedVersion != sha1ToDeploy {
	//		needToDeploy = true
	//	}
	//}

	//if needToDeploy == false {
	//	fmt.Fprintf(d.Writer, "Version %q is already deployed.", sha1ToDeploy)
	//	return
	//}

	fmt.Fprintf(d.Writer, "Going to deploy pod version %q \n", podVersion)
	for _, hook := range d.Context.Hooks {
		//@TODO Add a logger and log error coming from hooks
		err = hook.PreDeployment(userName, d.Env, d.Service, podVersion)
		if hook.StopOnError() == true && err != nil {
			fmt.Fprintf(d.Writer, "Predeployment failed: %q \n", err)
			return
		}
	}

	stateStream := make(chan deployers.State)
	go func() {
		for {
			state := <-stateStream
			//fmt.Fprintf(d.Writer, "%s : %d\n", state.Step, state.Progress)
			utils.RenderProgres(d.Writer, state.Step, state.Progress)

		}
	}()
	d.Context.Deployer.Deploy(d.Env, podVersion, stateStream)

	for _, hook := range d.Context.Hooks {
		//@TODO Add a logger and log error coming from hooks
		hook.PostDeployment(userName, d.Env, d.Service, podVersion)
	}

	fmt.Fprintln(d.Writer, "Deployment has successfully ended")
}

func (d *Deploy) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
	d.TableWriter = tablewriter.NewWriter(os.Stdout)
	d.Writer = os.Stdout
}
