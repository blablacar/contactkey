package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/blablacar/contactkey/context"
	"github.com/blablacar/contactkey/deployers"
	"github.com/blablacar/contactkey/utils"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	userName       = "Mister Robot"
	force          = false
	deployDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "contactkey_deploy_duration",
		Help: "The deploy duration of a project, in seconds",
	}, []string{"env", "project"})
)

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

func (d *Deploy) execute() error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(
		deployDuration.With(prometheus.Labels{"env": d.Env, "project": d.Service}).Set))
	d.Context.Metrics.Add(deployDuration)
	defer timer.ObserveDuration()

	currentUser, err := user.Current()
	if err == nil && currentUser.Name != "" {
		userName = currentUser.Name
	}

	if len(d.Context.PotentialUsername) > 0 {
		for _, value := range d.Context.PotentialUsername {
			potentialUserName := os.Getenv(value)
			if potentialUserName != "" {
				userName = potentialUserName
				break
			}
		}
	}

	if err := utils.CheckIfIsLaunchedInAScreen(); err != nil && d.Context.ScreenMandatory == true {
		return errors.New(fmt.Sprintf("Screen error raised: %q", err))
	}

	// The lock system is not mandatory
	if d.Context.LockSystem != nil {
		log.Println(fmt.Sprintf("Trying to lock the lock command for service %q and env %q", d.Service, d.Env))
		canLock, err := d.Context.LockSystem.Lock(d.Env, d.Service)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to lock, error raised: %q", err))
		}

		if canLock == false {
			return errors.New("Another command is currently running")
		}

		defer func(d *Deploy) {
			err = d.Context.LockSystem.Unlock(d.Env, d.Service)
			if err != nil {
				log.Errorln(fmt.Sprintf("Failed to unlock, error raised: %q", err))
			}
		}(d)
	}

	// If the branch is null it will use the default one.
	var podVersion string
	var sha1ToDeploy string
	if d.Context.GetVersionFromVcs {
		sha1ToDeploy, err = d.Context.Vcs.RetrieveSha1ForProject(branch)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to retrieve source changes for %q : %q", d.Service, err))
		}
		podVersion, err = d.Context.Binaries.RetrievePodVersion(sha1ToDeploy)
	} else {
		podVersion, err = d.Context.Binaries.RetrievePodVersion("")
	}
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to retrieve pod version: %q", err))
	}

	deployedVersions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to retrieve DEPLOYED version from service(%q) in env %q: %q", d.Service, d.Env, err))
	}

	if podVersion == "" {
		return errors.New(fmt.Sprintf("We have not found the pod version with the the sha1 %q \n"+
			"The pod has not been created.", sha1ToDeploy))
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
		return errors.New(fmt.Sprintf("Version %q is already deployed.", sha1ToDeploy))
	}

	log.Println(fmt.Sprintf("Going to deploy pod version %q \n", podVersion))
	for _, hook := range d.Context.Hooks {
		err = hook.PreDeployment(userName, d.Env, d.Service, podVersion)
		if hook.StopOnError() == true && err != nil {
			return errors.New(fmt.Sprintf("Predeployment failed: %q", err))
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
	err = d.Context.Deployer.Deploy(d.Env, podVersion, stateStream)
	if err != nil {
		return errors.New(fmt.Sprintf("Deployment failed: %q", err))
	}

	for _, hook := range d.Context.Hooks {
		err = hook.PostDeployment(userName, d.Env, d.Service, podVersion)
		if err != nil {
			log.Debugln(fmt.Sprintf("PostDeployment failed: %q", err))
		}
	}

	timer.ObserveDuration()
	err = d.Context.Metrics.Push()
	if err != nil {
		fmt.Fprintf(d.Writer, "Pushing metrics failed: %q \n", err)
	}

	log.Println(d.Writer, "Deployment has successfully ended")
	return nil
}

func (d *Deploy) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
	d.Writer = os.Stdout
}
