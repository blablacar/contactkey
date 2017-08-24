package deployers

import (
	"bufio"
	"encoding/json"
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
	Service      string
	Pod          string
	VcsRegexp    string
	WorkPath     string
	Environments map[string]string
	Log          *log.Logger
}

func NewDeployerGgn(cfg utils.DeployerGgnConfig,
	manifest utils.DeployerGgnManifest,
	logger *log.Logger) *DeployerGgn {
	return &DeployerGgn{
		WorkPath:     cfg.WorkPath,
		Service:      manifest.Service,
		Pod:          manifest.Pod,
		Environments: cfg.Environments,
		Log:          logger,
		VcsRegexp:    cfg.VcsRegexp,
	}
}

func (d *DeployerGgn) listUnits(env string) ([]string, error) {
	units := []string{}
	ggnCmd := ggn(env, "list-units")

	stdOut, err := ggnCmd.CombinedOutput()
	d.Log.WithFields(log.Fields{
		"cmd": ggnCmd.Path,
		"out": string(stdOut),
	}).Debug("Executing external command")

	if err != nil {
		d.Log.WithFields(log.Fields{
			"args": strings.Join(ggnCmd.Args, " "),
		}).Error("Failed to run external command")
		return nil, fmt.Errorf("Command `%s` failed with %q",
			strings.Join(ggnCmd.Args, " "),
			err)
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
		"cmd": strings.Join(ggnCmd.Args, " "),
		"out": string(stdOut),
	}).Debug("Executing external command")
	if err != nil {
		return "", err
	}
	return string(stdOut), nil
}

func (d *DeployerGgn) buildUnitRegexp(env string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s_%s_", env, d.Service))
}

func (d *DeployerGgn) buildVersionRegexp() *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`%s_aci-\S+:(\S+)`, d.Pod))
}

func (d *DeployerGgn) ListVersions(env string) (map[string]string, error) {
	ggnEnv, err := d.getGgnEnv(env)
	if err != nil {
		return nil, err
	}
	unitRegexp := d.buildUnitRegexp(ggnEnv)
	versionRegexp := d.buildVersionRegexp()
	versions := make(map[string]string)

	units, err := d.listUnits(ggnEnv)
	if err != nil {
		return nil, err
	}

	d.Log.Debugf("Matching for units with regex %q", unitRegexp)
	for _, unit := range units {
		if !unitRegexp.MatchString(unit) {
			continue
		}
		d.Log.Debugf("Found unit %q", unit)
		file, err := d.catUnit(ggnEnv, unit)
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
			d.Log.Debugf("Found version %q for unit %q", versions[unit], unit)
		}
	}

	return versions, nil
}

func (d *DeployerGgn) ListVcsVersions(env string) ([]string, error) {
	versions, err := d.ListVersions(env)
	if err != nil {
		return nil, err
	}
	regexp := regexp.MustCompile(d.VcsRegexp)
	vcsVersions := make([]string, 0)
	for _, version := range versions {
		vcsVersion := regexp.FindStringSubmatch(version)
		if len(vcsVersion) == 2 {
			vcsVersions = append(vcsVersions, vcsVersion[1])
		}
	}

	return vcsVersions, nil
}

func (d *DeployerGgn) getGgnEnv(env string) (string, error) {
	val, ok := d.Environments[env]
	if !ok {
		return "", fmt.Errorf("CCK env('%q') not found in GGN", env)
	}

	return val, nil
}

func (d *DeployerGgn) Deploy(env string, podVersion string, c chan State) error {
	serviceAttrs := make(map[string]string)

	ggnEnv, err := d.getGgnEnv(env)
	if err != nil {
		return err
	}

	serviceAttrs["version"] = podVersion
	serviceAttrsJSON, err := json.Marshal(serviceAttrs)
	if err != nil {
		return err
	}

	ggnCmd := ggn(ggnEnv, d.Service, "update", "-y", "-A", string(serviceAttrsJSON))
	d.Log.WithFields(log.Fields{
		"cmd": strings.Join(ggnCmd.Args, " "),
	})
	reader, err := utils.StreamCombinedOutput(ggnCmd)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	//statuses := States{}
	ggnCmd.Start()
	for scanner.Scan() {
		d.Log.Info(utils.VTClean(scanner.Text()))
		state := extractState(utils.VTClean(scanner.Text()))
		if state != (State{}) {
			c <- state
		}
	}

	ggnCmd.Wait()
	return nil
}

func extractState(ggnOutput string) State {
	s := State{}
	unitUpdate := regexp.MustCompile(`Remote service is already up to date .* unit=([\w\d]+)`)
	unitStart := regexp.MustCompile(`([\w\d]+) - Checking that unit is running.`)
	unitStartDone := regexp.MustCompile(`([\w\d]+): Ok - Deployed on`)

	if regexp.MustCompile(`Locking`).MatchString(ggnOutput) {
		s.Step = "locking"
		s.Progress = 100
	} else if m := unitUpdate.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] unit updating", m[1])
		s.Progress = 100
	} else if m := unitStart.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] unit start", m[1])
		s.Progress = 1
	} else if m := unitStartDone.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] unit start", m[1])
		s.Progress = 100
	} else if regexp.MustCompile(`Unlocking`).MatchString(ggnOutput) {
		s.Step = "unlocking"
		s.Progress = 100
	}
	return s
}

func (sts *States) updateStates(ggnOutput string) {
	unitUpdate := regexp.MustCompile(`Remote service is already up to date .* unit=([\w\d]+)`)
	unitStart := regexp.MustCompile(`([\w\d]+) - Checking that unit is running.`)
	unitStartDone := regexp.MustCompile(`([\w\d]+): Ok - Deployed on`)
	s := State{}

	if regexp.MustCompile(`Locking`).MatchString(ggnOutput) {
		s.Step = "locking"
		s.Progress = 100
	} else if m := unitUpdate.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] unit updating", m[1])
		s.Progress = 100
	} else if m := unitStart.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] unit start", m[1])
		s.Progress = 1
	} else if m := unitStartDone.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] unit start", m[1])
		s.Progress = 100
	} else if regexp.MustCompile(`Unlocking`).MatchString(ggnOutput) {
		s.Step = "unlocking"
		s.Progress = 100
	}

	if s != (State{}) {
		// update in place
		for i, is := range *sts {
			if is.Step == s.Step {
				(*sts)[i] = s
				return
			}
		}

		// or append
		*sts = append(*sts, s)
	}
}
