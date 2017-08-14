package utils

import (
	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

type Manifest struct {
	DeployerManifest          `mapstructure:"deployment"`
	VcsManifest               `mapstructure:"versionControlSystem"`
	RepositoryManagerManifest `mapstructure:"repositoryManager"`
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

type RepositoryManagerManifest struct {
	NexusManifest `mapstructure:"nexus"`
}

type NexusManifest struct {
	Artifact string `mapstructure:"artifact"`
}

type HookManifest struct {
	SlackManifest `mapstructure:"slack"`
}

type SlackManifest struct {
	Channel string `mapstructure:"channel"`
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
