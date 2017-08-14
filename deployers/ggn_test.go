package deployers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func mockggn(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestListUnits(t *testing.T) {
	execCommand = mockggn
	d := DeployerGgn{Log: log.New()}
	units, err := d.listUnits("staging")
	if err != nil {
		t.Fatal("listUnits failed")
	}
	if units[0] != "staging_webhooks_webhooks.service" {
		t.Errorf("Unexpected units[0] : %q", units[0])
	}
}

func TestCatUnit(t *testing.T) {
	execCommand = mockggn
	d := DeployerGgn{Log: log.New()}
	unit, err := d.catUnit("staging", "staging_webhooks_webhooks.service")
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
		Service: "wehooks",
	}
	r := d.buildUnitRegexp("staging")

	if r.String() != "^staging_wehooks_" {
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

func TestListVersions(t *testing.T) {
	envs := make(map[string]string)
	envs["staging"] = "staging"

	execCommand = mockggn
	d := DeployerGgn{Pod: "webhooks",
		Service:      "webhooks",
		Log:          log.New(),
		Environments: envs}
	v, err := d.ListVersions("staging")
	if err != nil {
		t.Fatal("listUnits failed")
	}
	if v["staging_webhooks_webhooks.service"] != "1.8.1-1" {
		t.Errorf("Unexpected unit version %q",
			v["staging_webhooks_webhooks.service"])
	}

}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := strings.Join(os.Args, " ")
	listUnits := regexp.MustCompile("-- ggn staging list-units$")
	catUnit := regexp.MustCompile("-- ggn staging fleetctl cat staging_webhooks_webhooks.service")
	switch {
	case listUnits.MatchString(args):
		fmt.Printf(`
staging_webhooks_webhooks.service						a34b757badea4abcaa518d0c686f82eb/10.13.35.193	active		running
staging_webhooks_webhooks.service					b102fa1e59ae42e28936dd676829236d/10.13.33.193	active		running
		`)
	case catUnit.MatchString(args):
		fmt.Printf(
			`ExecStart=/opt/bin/rkt      --insecure-options=all run      --set-env=TEMPLATER_OVERRIDE='${ATTR_0}'      --set-env=TEMPLATER_OVERRIDE_BASE64='${ATTR_BASE64_0}${ATTR_BASE64_1}'      --set-env=HOSTNAME='webhooks'      --set-env=HOST="%H"      --hostname=webhooks      --dns=10.254.0.3 --dns=10.254.0.4       --dns-search=pp-bourse.par-1.h.blbl.cr       --uuid-file-save=/mnt/sda9/rkt-uuid/pp-bourse/webhooks      --set-env=DOMAINNAME='pp.par-1.h.blbl.cr'      --net='bond0'      --set-env=AIRFLOW_HOME='/opt/webhooks'      aci.blbl.cr/pod-webhooks_aci-go-synapse:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-go-nerve:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-confd:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-embulk:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-zabbix-agent:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-webhooks:1.8.1-1 --exec /usr/local/bin/webhooks -- scheduler ---
	`)
	}
}

func TestListVcsVersions(t *testing.T) {
	envs := make(map[string]string)
	envs["staging"] = "staging"
	execCommand = mockggn
	d := DeployerGgn{PodName: "webhooks", Log: log.New(), Environments: envs, VcsRegexp: "-(.+)"}
	result, err := d.ListVcsVersions("staging")
	if err != nil {
		t.Fatalf("ListVcsVersions() failed : %q", err)
	}

	if len(result) != 1 {
		t.Fatal("We should have found only 1 version")
	}

	// We are receiving "1.8.1-1" when parsed with the regexp above it should "1" (string)
	if result[0] != "1" {
		t.Errorf("Result should be '1' instead found %q", result[0])
	}
}

func TestNewDeployerGGN(t *testing.T) {
}
