package context

import (
	"log"

	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/utils"
)

type Context struct {
	Deployer deployers.Deployer
	Log      log.Logger
}

func NewContext(cfg *utils.Config, manifest *utils.Manifest) *Context {
	ctx := &Context{}
	if manifest.DeployerManifest.DeployerGgnManifest != (utils.DeployerGgnManifest{}) {
		ctx.Deployer = deployers.NewDeployerGgn(
			cfg.DeployerConfig.DeployerGgnConfig,
			manifest.DeployerManifest.DeployerGgnManifest)
	}

	return ctx
}
