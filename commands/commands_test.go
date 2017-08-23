package commands

import (
	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/services"
	log "github.com/sirupsen/logrus"
)

type DeployerMockGgn struct {
	Log *log.Logger
}

func (d *DeployerMockGgn) ListVersions(env string) (map[string]string, error) {
	versions := map[string]string{
		"staging_webhooks_webhooks1.service": "26.1501244191-vb0f586a",
		"staging_webhooks_webhooks2.service": "26.1501244191-vb0f586a",
	}
	return versions, nil
}

func (d *DeployerMockGgn) ListVcsVersions(env string) ([]string, error) {
	vcsVersions := []string{"b0f586a", "b0f586a"}

	return vcsVersions, nil
}

func (d *DeployerMockGgn) Deploy(env string, podVersion string, c chan deployers.State) error {
	c <- deployers.State{
		Step:     "locking",
		Progress: 100,
	}

	return nil
}

type VCSMock struct {
}

func (v VCSMock) RetrieveSha1ForProject(branch string) (string, error) {
	return "vb0f586a", nil
}

func (v VCSMock) Diff(deployedSha1 string, sha1ToDeploy string) (*services.Changes, error) {
	return &services.Changes{}, nil
}

type RepositoryManagerMock struct{}

func (r RepositoryManagerMock) RetrievePodVersion(sha1 string) (string, error) {
	return "26.1501244191-vb0f586a", nil
}
