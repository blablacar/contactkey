package commands

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"io"

	"github.com/blablacar/contactkey/context"
	"github.com/blablacar/contactkey/deployers"
	"github.com/blablacar/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

func TestExecute(t *testing.T) {
	// Catch stdout
	out := new(bytes.Buffer)
	writer := io.Writer(out)
	log.SetOutput(writer)

	d := &Deploy{
		Context: &context.Context{
			Deployer: &DeployerMockGgn{},
			Vcs:      &SourcesMock{},
			Binaries: &BinariesMock{},
			Metrics:  utils.NewBlackholeMetricsRegistry(),
		},
		Writer: out,
	}
	d.execute()
	if out.String() == "" {
		t.Errorf("Unexpected stdout : %q", out)
	}

	if !regexp.MustCompile(`locking`).MatchString(out.String()) {
		t.Errorf("Stdout is missing locking step : %q", out)
	}

}

func TestUpdateSlow(t *testing.T) {
	ggnCmd := exec.Command("script", "-dp", "./testdata/already.script")
	reader, _ := utils.StreamCombinedOutput(ggnCmd)
	scanner := bufio.NewScanner(reader)
	ggnCmd.Start()

	for scanner.Scan() {
		state := deployers.ExtractState(utils.VTClean(scanner.Text()))
		if state != (deployers.State{}) {
			utils.RenderProgres(os.Stdout, state.Step, state.Progress)
		}
	}
	ggnCmd.Wait()
}
