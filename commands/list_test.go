package commands

import (
	"bytes"
	"testing"

	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
)

func TestListExecute(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := &List{
		Env:     "staging",
		Service: "webhooks",
		Context: &context.Context{
			Deployer: &DeployerMockGgn{},
		},
		Writer:      os.Stdout,
		TableWriter: tablewriter.NewWriter(out),
	}

	cmd.execute()
	if out.String() == "" {
		t.Errorf("Unexpected stdout : %q", out)
	}
}
