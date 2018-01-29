package context

import (
	"fmt"

	"io/ioutil"

	"github.com/blablacar/contactkey/deployers"
	"github.com/blablacar/contactkey/hooks"
	"github.com/blablacar/contactkey/services"
	"github.com/blablacar/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

type Context struct {
	Deployer          deployers.Deployer
	Vcs               services.Sources
	Binaries          services.Binaries
	Hooks             []hooks.Hooks
	LockSystem        utils.Lock
	ScreenMandatory   bool
	PotentialUsername []string
	Log               *log.Logger
	Metrics           *utils.MetricsRegistry
}

func NewContext(cfg *utils.Config, manifest *utils.Manifest) (*Context, error) {
	ctx := &Context{
		ScreenMandatory: cfg.ScreenMandatory,
	}

	log.SetLevel(log.DebugLevel)
	if cfg.Verbose == false {
		log.SetOutput(ioutil.Discard)
	}
	var err error
	if manifest.DeployerManifest.DeployerGgnManifest != (utils.DeployerGgnManifest{}) {
		ctx.Deployer, err = deployers.NewDeployerGgn(
			cfg.DeployerConfig.DeployerGgnConfig,
			manifest.DeployerManifest.DeployerGgnManifest)
		if err != nil {
			return nil, err
		}
	} else if manifest.DeployerManifest.DeployerK8sManifest != (utils.DeployerK8sManifest{}) {
		log.Debug("Creating a new DeployerK8s instance")
		ctx.Deployer, err = deployers.NewDeployerK8s(
			cfg.DeployerConfig.DeployerK8sConfig,
			manifest.DeployerManifest.DeployerK8sManifest)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(
			"Deployer unknown : %q",
			manifest.DeployerManifest,
		)
	}

	if manifest.VcsManifest.StashManifest != (utils.StashManifest{}) {
		ctx.Vcs, err = services.NewStash(
			cfg.VcsConfig.StashConfig,
			manifest.VcsManifest.StashManifest)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(
			"Vcs unknown : %q",
			manifest.VcsManifest,
		)
	}

	if manifest.BinariesManifest.NexusManifest != (utils.NexusManifest{}) {
		ctx.Binaries, err = services.NewNexus(
			cfg.Binaries.NexusConfig,
			manifest.BinariesManifest.NexusManifest)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Repository manager not found, You should check in your manifest if it's well formated.")
	}

	if manifest.HookManifest.SlackManifest != (utils.SlackManifest{}) {
		slack, err := hooks.NewSlack(
			cfg.HookConfig.SlackConfig,
			manifest.HookManifest.SlackManifest)

		if err != nil {
			return nil, err
		}

		ctx.Hooks = append(ctx.Hooks, slack)
	}
	if manifest.HookManifest.NewRelicManifest != (utils.NewRelicManifest{}) {
		newRelic, err := hooks.NewNewRelicClient(
			cfg.HookConfig.NewRelicConfig,
			manifest.HookManifest.NewRelicManifest)
		if err != nil {
			return nil, err
		}

		ctx.Hooks = append(ctx.Hooks, newRelic)
	}

	if len(manifest.HookManifest.ExecCommandManifest.OnPreDeploy) > 0 || len(manifest.HookManifest.ExecCommandManifest.OnPostDeploy) > 0 {
		ctx.Hooks = append(ctx.Hooks, hooks.NewExecommand(manifest.HookManifest.ExecCommandManifest))
	}

	if cfg.LockSystemConfig.FileLockConfig != (utils.FileLockConfig{}) {
		ctx.LockSystem, err = utils.NewFileLock(cfg.LockSystemConfig.FileLockConfig)

		if err != nil {
			return nil, err
		}
	}

	if len(cfg.PotentialUsername) > 0 {
		ctx.PotentialUsername = cfg.PotentialUsername
	}

	if cfg.MetricsConfig.PrometheusConfig != (utils.PrometheusConfig{}) {
		ctx.Metrics = utils.NewPrometheusMetricsRegistry(cfg.MetricsConfig.PrometheusConfig)
	} else {
		ctx.Metrics = utils.NewBlackholeMetricsRegistry()
	}
	return ctx, nil
}
