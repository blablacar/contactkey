package context

import (
	"fmt"
	"log"

	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/services"
	"github.com/remyLemeunier/contactkey/utils"
)

type Context struct {
	Deployer deployers.Deployer
	Vcs      services.VersionControlSystem
	Log      log.Logger
}

func NewContext(cfg *utils.Config, manifest *utils.Manifest) (*Context, error) {
	ctx := &Context{}
	if manifest.DeployerManifest.DeployerGgnManifest != (utils.DeployerGgnManifest{}) {
		ctx.Deployer = deployers.NewDeployerGgn(
			cfg.DeployerConfig.DeployerGgnConfig,
			manifest.DeployerManifest.DeployerGgnManifest)
	} else {
		return nil, fmt.Errorf(
			"Deployer unknown : %q",
			manifest.DeployerGgnManifest,
		)
	}

	if manifest.VcsManifest.StashManifest != (utils.StashManifest{}) {
		ctx.Vcs = services.NewStash(
			cfg.VcsConfig.StashConfig,
			manifest.VcsManifest.StashManifest)

	} else {
		return nil, fmt.Errorf(
			"Vcs unknown : %q",
			manifest.VcsManifest,
		)
	}

	return ctx, nil
}
