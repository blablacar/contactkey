package commands

import (
	"bytes"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/utils"
)

type DeployerMockGgn struct{}

func (d *DeployerMockGgn) ListVersions(env string) (map[string]string, error) {
	versions := map[string]string{
		"airflow1": "1",
		"airflow2": "1",
	}
	return versions, nil
}

func init() {
	deployers.Registry["mockggn"] = &DeployerMockGgn{}
}

func TestListExecute(t *testing.T) {
	cmd := &List{}
	out := new(bytes.Buffer)
	cmd.fill(
		&utils.Config{
			WorkPath: "./testdata/",
		},
		"airflow",
		"staging",
	)
	cmd.TableWriter = tablewriter.NewWriter(out)
	cmd.execute()
	if out.String() == "" {
		t.Errorf("Unexpected stdout : %q", out)
	}
}
