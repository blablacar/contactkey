package commands

import (
	"bytes"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	log "github.com/sirupsen/logrus"
)

type DeployerMockGgn struct {
	Log *log.Logger
}

func (d *DeployerMockGgn) Deploy(env string, podVersion string) error {
	return nil
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

func TestListExecute(t *testing.T) {
	cmd := &List{
		Env:     "staging",
		Service: "webhooks",
		Context: &context.Context{
			Deployer: &DeployerMockGgn{},
		},
	}
	out := new(bytes.Buffer)

	cmd.TableWriter = tablewriter.NewWriter(out)
	cmd.execute()
	if out.String() == "" {
		t.Errorf("Unexpected stdout : %q", out)
	}
}
