package deployers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func mockggn(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestDeployerGgnRegistry(t *testing.T) {
	if registry["github.com/remyLemeunier/contactkey/deployers.DeployerGgn"] == nil {
		t.Error("deployers.DeployerGgn is not in the registry")
	}
}

func TestListUnits(t *testing.T) {
	execCommand = mockggn
	d := DeployerGgn{}
	units, err := d.listUnits("staging")
	if err != nil {
		t.Fatal("listUnits failed")
	}
	if units[0] != "staging_airflow-flower_airflow-flower.service" {
		t.Error("Unexpected units[0] : %q", units[0])
	}
}

func TestCatUnit(t *testing.T) {
	execCommand = mockggn
	d := DeployerGgn{}
	unit, err := d.catUnit("staging", "staging_airflow-scheduler_airflow-scheduler.service")
	if err != nil {
		t.Fatal("CatUnits() failed : %q", err)
	}
	m, _ := regexp.MatchString("aci.blbl.cr/pod-airflow_aci-airflow:1.8.1-1", unit)
	if !m {
		t.Errorf("Unexpected CatUnit() : %q", unit)
	}

}

func TestListVersions(t *testing.T) {
	execCommand = mockggn
	d := DeployerGgn{name: "airflow-scheduler"}
	v, err := d.ListVersions("staging")
	if err != nil {
		t.Fatal("listUnits failed")
	}
	if v["staging_airflow-scheduler_airflow-scheduler.service"] != "1.8.1-1" {
		t.Errorf("Unexpected unit version %q",
			v["staging_airflow-scheduler_airflow-scheduler.service"])
	}

}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := strings.Join(os.Args, " ")

	listUnits := regexp.MustCompile("-- ggn staging list-units$")
	catUnit := regexp.MustCompile("-- ggn staging fleetctl cat staging_airflow-scheduler_airflow-scheduler.service")
	switch {
	case listUnits.MatchString(args):
		fmt.Printf(`
staging_airflow-flower_airflow-flower.service						a34b757badea4abcaa518d0c686f82eb/10.13.35.193	active		running
staging_airflow-scheduler_airflow-scheduler.service					b102fa1e59ae42e28936dd676829236d/10.13.33.193	active		running
		`)
	case catUnit.MatchString(args):
		fmt.Printf(`
ExecStart=/opt/bin/rkt      --insecure-options=all run      --set-env=TEMPLATER_OVERRIDE='${ATTR_0}'      --set-env=TEMPLATER_OVERRIDE_BASE64='${ATTR_BASE64_0}${ATTR_BASE64_1}'      --set-env=HOSTNAME='airflow-scheduler'      --set-env=HOST="%H"      --hostname=airflow-scheduler      --dns=10.254.0.3 --dns=10.254.0.4       --dns-search=pp-bourse.par-1.h.blbl.cr       --uuid-file-save=/mnt/sda9/rkt-uuid/pp-bourse/airflow-scheduler      --set-env=DOMAINNAME='pp.par-1.h.blbl.cr'      --net='bond0'      --set-env=AIRFLOW_HOME='/opt/airflow'      aci.blbl.cr/pod-airflow_aci-go-synapse:1.8.1-1    aci.blbl.cr/pod-airflow_aci-go-nerve:1.8.1-1    aci.blbl.cr/pod-airflow_aci-confd:1.8.1-1    aci.blbl.cr/pod-airflow_aci-embulk:1.8.1-1    aci.blbl.cr/pod-airflow_aci-zabbix-agent:1.8.1-1    aci.blbl.cr/pod-airflow_aci-airflow:1.8.1-1 --exec /usr/local/bin/airflow -- scheduler ---
	`)
	}
}
