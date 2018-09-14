package commands

import (
	"github.com/blablacar/contactkey/deployers"
	"github.com/blablacar/contactkey/services"
	log "github.com/sirupsen/logrus"
)

type DeployerMockGgn struct {
	Log *log.Logger
}

func (d *DeployerMockGgn) ListInstances(env string) ([]deployers.Instance, error) {
	return []deployers.Instance{
		{Name: "staging_webhooks_webhooks1.service", Version: "26.1501244191-vb0f586a"},
		{Name: "staging_webhooks_webhooks2.service", Version: "26.1501244191-vb0f586a"},
	}, nil
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

type SourcesMock struct {
}

func (v SourcesMock) RetrieveSha1ForProject(branch string) (string, error) {
	return "vb0f586a", nil
}

func (v SourcesMock) Diff(deployedSha1 string, sha1ToDeploy string, options services.DiffOptions) (*services.Changes, error) {
	return &services.Changes{}, nil
}

type BinariesMock struct{}

func (r BinariesMock) RetrievePodVersion(sha1 string) (string, error) {
	return "26.1501244191-vb0f586a", nil
}
