package context

import (
	"fmt"

	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/hooks"
	"github.com/remyLemeunier/contactkey/services"
	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

type Context struct {
	Deployer        deployers.Deployer
	Vcs             services.Sources
	Binaries        services.Binaries
	Hooks           []hooks.Hooks
	LockSystem      utils.Lock
	Log             *log.Logger
	ScreenMandatory bool
}

func NewContext(cfg *utils.Config, manifest *utils.Manifest) (*Context, error) {
	ctx := &Context{
		Log:             log.New(),
		ScreenMandatory: cfg.ScreenMandatory,
	}
	loglevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		loglevel = log.WarnLevel
	}
	ctx.Log.SetLevel(loglevel)

	if manifest.DeployerManifest.DeployerGgnManifest != (utils.DeployerGgnManifest{}) {
		ctx.Deployer, err = deployers.NewDeployerGgn(
			cfg.DeployerConfig.DeployerGgnConfig,
			manifest.DeployerManifest.DeployerGgnManifest,
			ctx.Log)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(
			"Deployer unknown : %q",
			manifest.DeployerGgnManifest,
		)
	}

	if manifest.VcsManifest.StashManifest != (utils.StashManifest{}) {
		ctx.Vcs, err = services.NewStash(
			cfg.VcsConfig.StashConfig,
			manifest.VcsManifest.StashManifest,
			ctx.Log)
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
			manifest.BinariesManifest.NexusManifest,
			ctx.Log)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Repository manager not found, You should check in your manifest if it's well formated.")
	}

	if manifest.HookManifest.SlackManifest != (utils.SlackManifest{}) {
		slack, err := hooks.NewSlack(
			cfg.HookConfig.SlackConfig,
			manifest.HookManifest.SlackManifest,
			ctx.Log)

		if err != nil {
			return nil, err
		}

		ctx.Hooks = append(ctx.Hooks, slack)
	}
	if manifest.HookManifest.NewRelicManifest != (utils.NewRelicManifest{}) {
		newRelic, err := hooks.NewNewRelicClient(
			cfg.HookConfig.NewRelicConfig,
			manifest.HookManifest.NewRelicManifest,
			ctx.Log)

		if err != nil {
			return nil, err
		}

		ctx.Hooks = append(ctx.Hooks, newRelic)
	}

	if len(manifest.HookManifest.ExecCommandManifest.OnPreDeploy) > 0 || len(manifest.HookManifest.ExecCommandManifest.OnPostDeploy) > 0 {
		ctx.Hooks = append(ctx.Hooks, hooks.NewExecommand(
			manifest.HookManifest.ExecCommandManifest,
			ctx.Log))
	}

	if cfg.LockSystemConfig.FileLockConfig != (utils.FileLockConfig{}) {
		ctx.LockSystem, err = utils.NewFileLock(cfg.LockSystemConfig.FileLockConfig)

		if err != nil {
			return nil, err
		}
	}

	return ctx, nil
}
