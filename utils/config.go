package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

var DefaultHome = filepath.Join(os.Getenv("HOME"), ".contactkey", "config.yml")

type Config struct {
	WorkPath           string `mapstructure:"workPath"`
	LogLevel           string
	GlobalEnvironments []string
	DeployerConfig     `mapstructure:"deployers"`
	VcsConfig          `mapstructure:"versionControlSystem"`
	RepositoryManager  `mapstructure:"repositoryManager"`
	HookConfig         `mapstructure:"hooks"`
}

type DeployerConfig struct {
	DeployerGgnConfig `mapstructure:"ggn"`
}

type DeployerGgnConfig struct {
	WorkPath     string            `mapstructure:"workPath"`
	Environments map[string]string `mapstructure:"environments"`
	VcsRegexp    string            `mapstructure:"vcsRegexp"`
}

type VcsConfig struct {
	StashConfig `mapstructure:"stash"`
}

type StashConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Url      string `mapstructure:"url"`
}

type RepositoryManager struct {
	NexusConfig `mapstructure:"nexus"`
}

type NexusConfig struct {
	Url           string `mapstructure:"url"`
	Repository    string `mapstructure:"repository"`
	Group         string `mapstructure:"group"`
	ServiceRegexp string `mapstructure:"serviceRegexp"`
}

type HookConfig struct {
	SlackConfig `mapstructure:"slack"`
}

type SlackConfig struct {
	Url   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

func LoadConfig(cfgReader []byte) (*Config, error) {
	cfg := &Config{}
	cfgAux := make(map[string]interface{})
	err := yaml.Unmarshal(cfgReader, &cfgAux)
	if err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(cfgAux, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c Config) DiscoverServices() ([]string, error) {
	services := make([]string, 0)

	files, err := ioutil.ReadDir(c.WorkPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() == true {
			continue
		}

		ext := filepath.Ext(file.Name())
		if ext == ".yaml" || ext == ".yml" {
			baseNameWithoutExt := strings.TrimSuffix(filepath.Base(file.Name()), ext)
			services = append(services, baseNameWithoutExt)
		}
	}

	return services, nil
}

func ReadFile(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}
