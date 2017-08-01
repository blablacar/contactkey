package utils

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	WorkPath           string
	GlobalEnvironments []string
	Deployers          ConfigDeployers `mapstructure:"deployers"`
	Environment        interface{}     `mapstructure:"environments"`
	DeployDefaults     DeployManifest  `yaml:"deployDefaults"`
}

type ConfigDeployers struct {
	DeployerGgn ConfigDeployerGgn `mapstructure:"ggn"`
}

type ConfigDeployerGgn struct {
	GitBuildtoolsUrl     string            `yaml:"gitBuildtoolsUrl"`
	WorkPath             string            `yaml:"workPath,omitempty"`
	User                 string            `yaml:"user,omitempty"`
	SupportedEnvironment map[string]string `yaml:"supportedEnvironment"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := new(Config)

	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".contactkey"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
