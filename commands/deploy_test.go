package commands

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/remyLemeunier/contactkey/context"
	"github.com/remyLemeunier/contactkey/deployers"
	"github.com/remyLemeunier/contactkey/utils"
)

func TestExecute(t *testing.T) {
	out := new(bytes.Buffer)
	d := &Deploy{
		Context: &context.Context{
			Deployer:          &DeployerMockGgn{},
			Vcs:               &SourcesMock{},
			Binaries: &BinariesMock{},
		},
		Writer: out,
	}
	d.execute()
	if out.String() == "" {
		t.Errorf("Unexpected stdout : %q", out)
	}

	if !regexp.MustCompile(`locking : 100`).MatchString(out.String()) {
		t.Errorf("Stdout is missing locking step : %q", out)
	}

}

func TestUpdateSlow(t *testing.T) {
	//sts := deployers.States{}

	ggnCmd := exec.Command("script", "-dp", "./testdata/already.script")
	reader, _ := utils.StreamCombinedOutput(ggnCmd)
	scanner := bufio.NewScanner(reader)
	ggnCmd.Start()

	for scanner.Scan() {
		state := deployers.ExtractState(utils.VTClean(scanner.Text()))
		if state != (deployers.State{}) {
			//fmt.Printf("%s : %d\n", state.Step, state.Progress)
			utils.RenderProgres(os.Stdout, state.Step, state.Progress)
		}
	}
	ggnCmd.Wait()
}
