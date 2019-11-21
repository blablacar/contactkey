package deployers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"errors"

	"github.com/blablacar/contactkey/utils"
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
}

type FleetUnit struct {
	name   string
	active string
	sub    string
}

func NewDeployerGgn(cfg utils.DeployerGgnConfig, manifest utils.DeployerGgnManifest) (*DeployerGgn, error) {
	if manifest.Service == "" {
		return nil, errors.New("You need to define a service name for ggn in the manifest.")
	}

	if manifest.Pod == "" {
		return nil, errors.New("You need to define a pod name for ggn in the manifest.")
	}

	if len(cfg.Environments) == 0 {
		return nil, errors.New("You need to define at least a pair of env for ggn in the config(Array between cck env and ggn env).")
	}

	return &DeployerGgn{
		WorkPath:     cfg.WorkPath,
		Service:      manifest.Service,
		Pod:          manifest.Pod,
		Environments: cfg.Environments,
		VcsRegexp:    cfg.VcsRegexp,
	}, nil
}

func (d *DeployerGgn) catUnit(env string, unit string) (string, error) {
	ggnCmd := ggn(env, "fleetctl", "cat", unit)
	stdOut, err := ggnCmd.CombinedOutput()
	log.WithFields(log.Fields{
		"cmd": strings.Join(ggnCmd.Args, " "),
		"out": string(stdOut),
	}).Debug("Executing external command")
	if err != nil {
		return "", err
	}
	return string(stdOut), nil
}

func (d *DeployerGgn) buildUnitRegexp(env string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`%s_%s_`, env, d.Service))
}

func (d *DeployerGgn) buildVersionRegexp() *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`%s_aci-\S+:(\S+)`, d.Pod))
}

func (d *DeployerGgn) InstanceState() string {
	return ""
}

func (d *DeployerGgn) ListInstances(env string) ([]Instance, error) {
	instances := []Instance{}
	ggnEnv, err := d.getGgnEnv(env)
	if err != nil {
		return nil, err
	}
	unitRegexp := d.buildUnitRegexp(ggnEnv)

	units, err := d.listFleetUnits(ggnEnv)
	if err != nil {
		return nil, err
	}

	log.Debugf("Matching units against regex %q", unitRegexp)
	for _, unit := range units {
		if !unitRegexp.MatchString(unit.name) {
			continue
		}
		log.Debugf("Found unit %q", unit.name)
		file, err := d.catUnit(ggnEnv, unit.name)
		if err != nil {
			log.Warnf("Failed to cat unit %v", unit.name)
			continue
		}
		if m := d.buildVersionRegexp().FindStringSubmatch(file); m != nil {
			instances = append(instances, Instance{Name: unit.name, State: unit.sub, Version: m[1]})
			log.Debugf("Found instance %q", unit.name)
		}
	}

	return instances, nil
}

func (d *DeployerGgn) ListVcsVersions(env string) ([]string, error) {
	instances, err := d.ListInstances(env)
	if err != nil {
		return nil, err
	}
	regexp := regexp.MustCompile(d.VcsRegexp)
	vcsVersions := make([]string, 0)
	for _, instance := range instances {
		vcsVersion := regexp.FindStringSubmatch(instance.Version)
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
	log.WithFields(log.Fields{
		"cmd": strings.Join(ggnCmd.Args, " "),
	})
	reader, err := utils.StreamCombinedOutput(ggnCmd)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	ggnCmd.Start()
	for scanner.Scan() {
		log.Debug(utils.VTClean(scanner.Text()))
		state := ExtractState(utils.VTClean(scanner.Text()))
		if state != (State{}) {
			c <- state
		}
	}

	ggnCmd.Wait()
	return nil
}

func ExtractState(ggnOutput string) State {
	s := State{}
	unitUpdate := regexp.MustCompile(`Remote service is already up to date .* unit=([\w\d]+)`)
	unitStart := regexp.MustCompile(`([\w\d]+) - Checking that unit is running.`)
	unitStartDone := regexp.MustCompile(`([\w\d]+): Ok - Deployed on`)
	unitHealthy := regexp.MustCompile(`\[ZkCheck\].* ([\w\d]+) .* - Checking that service adds key in zookeeper`)
	unitHealthyDone := regexp.MustCompile(`\[ZkCheck\].* ([\w\d]+) - .*: Ok`)

	if regexp.MustCompile(`Locking`).MatchString(ggnOutput) {
		s.Step = "locking deployment"
		s.Progress = 100
	} else if m := unitUpdate.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] updating service in fleet", m[1])
		s.Progress = 100
	} else if m := unitStart.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] starting instance", m[1])
		s.Progress = 1
	} else if m := unitStartDone.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] starting instance", m[1])
		s.Progress = 100
	} else if m := unitHealthy.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] checking instance health", m[1])
		s.Progress = 1
	} else if m := unitHealthyDone.FindStringSubmatch(ggnOutput); m != nil {
		s.Step = fmt.Sprintf("[%s] checking instance health", m[1])
		s.Progress = 100
	} else if regexp.MustCompile(`Unlocking`).MatchString(ggnOutput) {
		s.Step = "unlocking deployment"
		s.Progress = 100
	} else if regexp.MustCompile(`Victory !`).MatchString(ggnOutput) {
		s.Step = "victory"
		s.Progress = 100
	}

	return s
}

func (d DeployerGgn) listFleetUnits(ggnEnv string) ([]FleetUnit, error) {
	units := []FleetUnit{}
	unitRegexp := regexp.MustCompile(`(\S+)\t+(\S+)\t+(\S+)\t+(\S+)\t+(\S+)\t+(\S+)`)
	cmd := ggn(ggnEnv, "fleetctl", "--", "list-units", "--no-legend", "--fields='active,hash,load,machine,sub,unit'")

	log.WithFields(log.Fields{"cmd": strings.Join(cmd.Args, " ")}).Debug("Executing external command")

	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(reader)
	cmd.Start()

	for scanner.Scan() {
		if m := unitRegexp.FindStringSubmatch(scanner.Text()); m != nil {
			units = append(units, FleetUnit{
				active: m[1],
				sub:    m[5],
				name:   m[6],
			})
		}

	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return units, nil
}
