package context

import (
	"fmt"
	"log"

	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/utils"
)

type Context struct {
	Deployer deployers.Deployer
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

	return ctx, nil
}
