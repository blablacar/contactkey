package deployers

import (
	"errors"
	"fmt"
	"path"
	"regexp"

	"github.com/remyLemeunier/contactkey/utils"
	"github.com/remyLemeunier/k8s-deploy/releases"

	log "github.com/sirupsen/logrus"
)

type DeployerK8s struct {
	ReleaseName  string
	Namespace    string
	Release      releases.Releaser
	chartPath    string
	valueFiles   []string
	values       []string
	Environments map[string]utils.K8sEnvironment
	workPath     string
	vcsRegexp    string
}

func NewDeployerK8s(cfg utils.DeployerK8sConfig, manifest utils.DeployerK8sManifest) (*DeployerK8s, error) {
	var deployer DeployerK8s

	if manifest.Release == "" {
		return nil, errors.New("You need to define a release name for k8s")
	}

	deployer.ReleaseName = manifest.Release
	deployer.Namespace = manifest.Namespace
	deployer.Environments = cfg.Environments
	deployer.workPath = cfg.WorkPath
	deployer.vcsRegexp = cfg.VcsRegexp
	return &deployer, nil
}

func (d *DeployerK8s) ListInstances(env string) ([]Instance, error) {
	context, err := d.getK8sContext(env)
	if err != nil {
		return nil, err
	}
	log.Debug("ListInstances")
	log.Debugf("Environment %s matched to : cluster=%s", env, context.Cluster)

	if d.Release == nil {
		d.Release, err = d.getRelease(context.Cluster, d.Namespace)
		if err != nil {
			return nil, err
		}
	}

	status, err := d.Release.Status()
	if err != nil {
		return nil, err
	}

	content, err := d.Release.Content()
	if err != nil {
		return nil, err
	}

	return []Instance{
		{
			Name:    status.Name,
			State:   status.Info.Status.Code.String(),
			Version: content.Release.Chart.Metadata.Version,
		},
	}, nil
}

func (d *DeployerK8s) ListVcsVersions(env string) ([]string, error) {
	context, err := d.getK8sContext(env)
	if err != nil {
		return nil, err
	}

	log.Debug("ListVcsVersions")
	log.Debugf("Environment %s matched to : cluster=%s", env, context.Cluster)

	if d.Release == nil {
		d.Release, err = d.getRelease(context.Cluster, d.Namespace)
		if err != nil {
			return nil, err
		}
	}

	content, err := d.Release.Content()
	if err != nil {
		return nil, err
	}

	vcsVersion := regexp.MustCompile(d.vcsRegexp).FindStringSubmatch(content.String())
	if len(vcsVersion) != 2 {
		return nil, fmt.Errorf("Could not find vcsVersion")
	}

	return []string{vcsVersion[1]}, nil
}

func (d *DeployerK8s) getK8sContext(env string) (utils.K8sEnvironment, error) {
	if d.Environments[env] == (utils.K8sEnvironment{}) {
		return utils.K8sEnvironment{}, fmt.Errorf("Cannot find K8s context for environment : %s", env)
	}
	return d.Environments[env], nil
}

func (d *DeployerK8s) getReleasePath(cluster string, namespace string) string {
	return path.Join(d.workPath, "deployments", cluster, namespace, d.ReleaseName, "release.yaml")

}

func (d *DeployerK8s) getRelease(cluster string, namespace string) (*releases.Release, error) {
	manifest := d.getReleasePath(cluster, namespace)
	return releases.NewReleaseFromManifest(manifest)
}

func (d *DeployerK8s) Deploy(env string, podVersion string, c chan State) error {
	context, err := d.getK8sContext(env)
	if err != nil {
		return err
	}
	log.Debugf("Environment %s matched to : cluster=%s", env, context.Cluster)

	if d.Release == nil {
		d.Release, err = d.getRelease(context.Cluster, d.Namespace)
		if err != nil {
			return err
		}
	}

	overrides := []string{
		fmt.Sprintf("image.tag=%s", podVersion),
	}
	d.Release.AddValues([]string{}, overrides)

	log.Debugf("Deploying %s, version %s", d.ReleaseName, podVersion)
	c <- State{Step: "deploying release", Progress: 0}
	err = d.Release.Deploy()
	if err != nil {
		return err
	}
	c <- State{Step: "deploying release", Progress: 100}
	return nil
}
