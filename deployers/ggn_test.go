package deployers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/blablacar/contactkey/deployers/testdata"
)

func mockggn(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := strings.Join(os.Args, " ")
	switch {
	case testdata.CatUnitRegexp.MatchString(args):
		fmt.Printf(testdata.CatUnit)
	case testdata.ListFleetUnitsRegexp.MatchString(args):
		fmt.Printf(testdata.ListFleetUnits)
	}
}

func TestCatUnit(t *testing.T) {
	execCommand = mockggn
	d := DeployerGgn{}
	unit, err := d.catUnit("staging", "staging_webhooks_webhooks1.service")
	if err != nil {
		t.Fatalf("CatUnits() failed : %q", err)
	}
	m, _ := regexp.MatchString("aci.blbl.cr/pod-webhooks_aci-webhooks:1.8.1-1", unit)
	if !m {
		t.Errorf("Unexpected CatUnit() : %q", unit)
	}

}

func TestBuildUnitRegexp(t *testing.T) {
	d := DeployerGgn{
		Service: "webhooks",
	}
	r := d.buildUnitRegexp("staging")

	if r.String() != "staging_webhooks_" {
		t.Errorf("Unexpected UnitRegexp : %q", r)
	}
}

func TestVersionRegexp(t *testing.T) {
	d := DeployerGgn{
		Pod: "aci.blbl.cr/pod-webhooks",
	}
	r := d.buildVersionRegexp()

	if r.String() != "aci.blbl.cr/pod-webhooks_aci-\\S+:(\\S+)" {
		t.Errorf("Unexpected UnitRegexp : %q", r)
	}
}

func TestListInstances(t *testing.T) {
	envs := make(map[string]string)
	envs["staging"] = "staging"

	execCommand = mockggn
	d := DeployerGgn{Pod: "webhooks",
		Service:      "webhooks",
		Environments: envs}
	i, err := d.ListInstances("staging")
	if err != nil {
		t.Fatal("listUnits failed")
	}
	if i[1].Version != "1.8.1-1" {
		t.Errorf("Unexpected unit version %q", i[1])
	}
}

func TestListVcsVersions(t *testing.T) {
	envs := make(map[string]string)
	envs["staging"] = "staging"
	execCommand = mockggn
	d := DeployerGgn{
		Service:      "webhooks",
		Pod:          "pod-webhooks",
		Environments: envs,
		VcsRegexp:    "-(.+)",
	}
	result, err := d.ListVcsVersions("staging")
	if err != nil {
		t.Fatalf("ListVcsVersions() failed : %q", err)
	}

	if len(result) != 3 {
		t.Errorf("We should have found only 3 versions, found %q", result)
	}

	// We are receiving "1.8.1-1" when parsed with the regexp above it should "1" (string)
	if result[0] != "1" {
		t.Errorf("Result should be '1' instead found %q", result[0])
	}
}

func TestExtractState(t *testing.T) {
	s := State{}

	s = ExtractState("[ZkCheck][webhooks] webhooks2 /services/wehooks - Checking that service adds key in zookeeper")
	if s.Step != "[webhooks2] checking instance health" {
		t.Errorf("Unexpected step : %q", s.Step)
	}
	s = ExtractState("[ZkCheck][webhooks] webhooks2 - /services/webhooks: Ok")
}

func TestListFleetUnits(t *testing.T) {
	d := DeployerGgn{}
	units, err := d.listFleetUnits("staging")
	if err != nil {
		t.Fatalf("Unexpected error : %q", err)
	}
	if len(units) != 4 {
		t.Errorf("Unexpected units : %q", units)
	}
	if units[0].name != "staging_sleepy_sleepy0.service" {
		t.Errorf("Unexpected unit[0] : %q", units[0])
	}
}
