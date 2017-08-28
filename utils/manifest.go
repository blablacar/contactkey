package utils

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

type Manifest struct {
	DeployerManifest          `mapstructure:"deployment"`
	VcsManifest               `mapstructure:"sources"`
	BinariesManifest `mapstructure:"binaries"`
	HookManifest              `mapstructure:"hooks"`
}

type VcsManifest struct {
	StashManifest `mapstructure:"stash"`
}

type StashManifest struct {
	Repository string `mapstructure:"repository"`
	Project    string `mapstructure:"project"`
	Branch     string `mapstructure:"branch"`
}

type DeployerManifest struct {
	DeployerGgnManifest  `mapstructure:"ggn"`
	DeployerHelmManifest `mapstructure:"helm"`
}

type DeployerGgnManifest struct {
	Service string `mapstructure:"service"`
	Pod     string `mapstructure:"pod"`
}

type DeployerHelmManifest struct {
	ReleaseName string `mapstructure:"release"`
}

type BinariesManifest struct {
	NexusManifest `mapstructure:"nexus"`
}

type NexusManifest struct {
	Artifact string `mapstructure:"artifact"`
}

type HookManifest struct {
	SlackManifest       `mapstructure:"slack"`
	NewRelicManifest    `mapstructure:"newRelic"`
	ExecCommandManifest `mapstructure:"execCommand"`
}

type SlackManifest struct {
	Channel     string `mapstructure:"channel"`
	StopOnError bool   `mapstructure:"stopOnError"`
}

type NewRelicManifest struct {
	ApplicationId string `mapstructure:"applicationId"`
	StopOnError   bool   `mapstructure:"stopOnError"`
}

type ExecCommandManifest struct {
	OnPreDeploy  []CommandList `mapstructure:"onPreDeploy"`
	OnPostDeploy []CommandList `mapstructure:"onPostDeploy"`
	StopOnError  bool          `mapstructure:"stopOnError"`
}

type CommandList struct {
	Command string   `mapstructure:"command"`
	Args    []string `mapstructure:"args"`
}

func LoadManifest(manifestReader []byte) (*Manifest, error) {
	manifest := &Manifest{}
	manifestAux := make(map[string]interface{})
	err := yaml.Unmarshal(manifestReader, &manifestAux)
	if err != nil {
		return nil, err
	}

	if err := mapstructure.Decode(manifestAux, manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

func (m *DeployerGgnManifest) validate() error {
	if m.Service == "" {
		return fmt.Errorf("Missing field Service")
	}
	return nil
}
