package deployers

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

var execCommand = exec.Command

func ggn(args ...string) *exec.Cmd {
	return execCommand("ggn", args...)
}

type DeployerGgn struct {
	PodName  string
	WorkPath string
	Log      *log.Logger
}

func NewDeployerGgn(cfg utils.DeployerGgnConfig, manifest utils.DeployerGgnManifest) *DeployerGgn {
	return &DeployerGgn{
		WorkPath: cfg.WorkPath,
		PodName:  manifest.PodName,
		Log:      log.New(),
	}
}

func init() {
	Registry["ggn"] = &DeployerGgn{Log: log.New()}
}

func (d *DeployerGgn) SetLogLevel(level log.Level) {
	d.Log.SetLevel(level)
}

func (d *DeployerGgn) listUnits(env string) ([]string, error) {
	units := []string{}
	ggnCmd := ggn(env, "list-units")

	stdOut, err := ggnCmd.CombinedOutput()
	d.Log.WithFields(log.Fields{
		"cmd":    ggnCmd.Path,
		"args":   strings.Join(ggnCmd.Args, " "),
		"stdout": ggnCmd.Stdout,
	}).Debug("Executing external command")
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(".*.service")
	for _, line := range strings.Split(string(stdOut), "\n") {
		unit := r.FindString(line)
		if unit != "" {
			units = append(units, unit)
		}
	}

	return units, nil

}

func (d *DeployerGgn) catUnit(env string, unit string) (string, error) {
	ggnCmd := ggn(env, "fleetctl", "cat", unit)
	stdOut, err := ggnCmd.CombinedOutput()
	d.Log.WithFields(log.Fields{
		"cmd":    ggnCmd.Path,
		"args":   strings.Join(ggnCmd.Args, " "),
		"stdout": ggnCmd.Stdout,
	}).Debug("Executing external command")
	if err != nil {
		return "", err
	}
	return string(stdOut), nil
}

func (d *DeployerGgn) ListVersions(env string) (map[string]string, error) {
	unitRegexp := regexp.MustCompile(fmt.Sprintf("_%s_", "webhooks"))
	versionRegexp := regexp.MustCompile("pod-webhooks_aci-\\S+:(\\S+)")
	versions := make(map[string]string)

	units, err := d.listUnits(env)
	if err != nil {
		return nil, err
	}
	for _, unit := range units {
		d.Log.Debug(unit)
		if !unitRegexp.MatchString(unit) {
			continue
		}
		file, err := d.catUnit(env, unit)
		if err != nil {
			// @TODO what should we do there ?
			continue
		}
		version := versionRegexp.FindStringSubmatch(file)
		if len(version) == 0 {
			continue
		}
		if version[1] != "" {
			versions[unit] = version[1]
		}
	}

	return versions, nil
}

func (d *DeployerGgn) Deploy(env string) error {
	// @TODO

	return nil
}
