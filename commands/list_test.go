package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	log "github.com/sirupsen/logrus"
)

func TestListExecute(t *testing.T) {
	// Catch stdout
	out := new(bytes.Buffer)
	logger := log.New()
	logger.Out = out

	cmd := &List{
		Env:     "staging",
		Service: "webhooks",
		Context: &context.Context{
			Deployer: &DeployerMockGgn{},
			Log:      logger,
		},
		Writer:      os.Stdout,
		TableWriter: tablewriter.NewWriter(out),
	}

	cmd.execute()
	if out.String() == "" {
		t.Errorf("Unexpected stdout : %q", out)
	}
}
