package commands

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	log "github.com/sirupsen/logrus"
)

func TestListExecute(t *testing.T) {
	// Catch stdout
	out := new(bytes.Buffer)
	writer := io.Writer(out)
	log.SetOutput(writer)

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
