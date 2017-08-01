package utils

import (
	"io/ioutil"

	"github.com/imdario/mergo"

	"gopkg.in/yaml.v2"
)

const ManifestVersion = "v1"

type DeployManifest struct {
	ManifestVersion string      `yaml:"manifestVersion"`
	Stash           interface{} `yaml:"stash"`
	Deploy          Deploy      `yaml:"deployment" mapstructure:"deployment"`
}

type Stash struct {
	Project       string `yaml:"project"`
	Repo          string `yaml:"repo"`
	DefaultBranch string `yaml:"defaultBranch"`
}

type Deploy struct {
	Method      string      `yaml:"method"`
	ServiceName string      `yaml:"serviceName"`
	PodName     string      `yaml:"podName"`
	Hooks       DeployHooks `yaml:"hooks"`
}

type DeployHooks struct {
	DeployHookSlack    `yaml:"slack,omitempty" mapstructure:"slack"`
	DeployHookNewRelic `yaml:"newRelic,omitempty" mapstructure:"newRelic"`
}

type DeployHookNewRelic struct {
	ApiKey  string `yaml:"apiKey"`
	AppName string `yaml:"admin-tools"`
}

type DeployHookSlack struct {
	Channels []string `yaml:"channels"`
}

func LoadDeployfile(defaults *DeployManifest, filename string) (*DeployManifest, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	manifest, err := UnmarshalDeployfile(b)
	if err != nil {
		return nil, err
	}

	err = mergo.MergeWithOverwrite(manifest, defaults)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func UnmarshalDeployfile(data []byte) (*DeployManifest, error) {
	y := &DeployManifest{}
	err := yaml.Unmarshal(data, y)
	if err != nil {
		return nil, err
	}
	return y, nil
}
