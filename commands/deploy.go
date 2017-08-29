package commands

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/remyLemeunier/contactkey/context"
	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var userName = "Mister Robot"
var force = false

func init() {
	deployCmd.PersistentFlags().StringVar(&branch, "branch", "", "Change the branch from the default one.")
	deployCmd.PersistentFlags().BoolVar(&force, "force", false, "Force the deployement, even if already up to date")
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the service in an environment",
}

type Deploy struct {
	Env     string
	Service string
	Context *context.Context
	Writer  io.Writer
}

func (d *Deploy) execute() {
	currentUser, err := user.Current()
	if err == nil {
		userName = currentUser.Name
	}

	if err := utils.CheckIfIsLaunchedInAScreen(); err != nil && d.Context.ScreenMandatory == true {
		log.Errorln(fmt.Sprintf("Screen error raised: %q", err))
		return
	}

	// The lock system is not mandatory
	if d.Context.LockSystem != nil {
		log.Println(fmt.Sprintf("Trying to lock the lock command for service %q and env %q", d.Service, d.Env))
		canLock, err := d.Context.LockSystem.Lock(d.Env, d.Service)
		if err != nil {
			log.Errorln(fmt.Sprintf("Failed to lock, error raised: %q", err))
			return
		}

		if canLock == false {
			log.Errorln("Another command is currently running")
			return
		}

		defer func(d *Deploy) {
			d.Context.LockSystem.Unlock(d.Env, d.Service)
			if err != nil {
				log.Errorln(fmt.Sprintf("Failed to unlock, error raised: %q", err))
			}
		}(d)
	}

	// If the branch is null it will use the default one.
	sha1ToDeploy, err := d.Context.Vcs.RetrieveSha1ForProject(branch)
	if err != nil {
		log.Errorln(fmt.Sprintf("Failed to retrieve source changes for %q : %q", d.Service, err))
		return
	}

	podVersion, err := d.Context.Binaries.RetrievePodVersion(sha1ToDeploy)
	if err != nil {
		log.Errorln(fmt.Sprintf("Failed to retrieve pod version: %q", err))
		return
	}

	deployedVersions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	if err != nil {
		log.Errorln(fmt.Sprintf("Failed to retrieve DEPLOYED version from service(%q) in env %q: %q", d.Service, d.Env, err))
		return
	}

	if podVersion == "" {
		log.Errorln(fmt.Sprintf("We have not found the pod version with the the sha1 %q \n"+
			"The pod has not been created.", sha1ToDeploy))
		return
	}

	needToDeploy := force
	for _, deployedVersion := range deployedVersions {
		if deployedVersion != sha1ToDeploy {
			needToDeploy = true
		}
	}

	if len(deployedVersions) == 0 {
		needToDeploy = true
	}

	if needToDeploy == false {
		log.Errorln(fmt.Sprintf("Version %q is already deployed.", sha1ToDeploy))
		return
	}

	log.Println(fmt.Sprintf("Going to deploy pod version %q \n", podVersion))
	for _, hook := range d.Context.Hooks {
		err = hook.PreDeployment(userName, d.Env, d.Service, podVersion)
		if hook.StopOnError() == true && err != nil {
			log.Errorln(fmt.Sprintf("Predeployment failed: %q", err))
			return
		} else if err != nil {
			log.Debugln(fmt.Sprintf("Predeployment failed: %q", err))
		}
	}

	stateStream := make(chan deployers.State)
	go func() {
		for {
			state := <-stateStream
			utils.RenderProgres(d.Writer, state.Step, state.Progress)

		}
	}()
	d.Context.Deployer.Deploy(d.Env, podVersion, stateStream)

	for _, hook := range d.Context.Hooks {
		err = hook.PostDeployment(userName, d.Env, d.Service, podVersion)
		if err != nil {
			log.Debugln(fmt.Sprintf("PostDeployment failed: %q", err))
		}
	}

	log.Println("Deployment has successfully ended")
}

func (d *Deploy) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
	d.Writer = os.Stdout
}
