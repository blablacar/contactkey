package utils

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	Deployers      ConfigDeployers `mapstructure:"deployers"`
	Environment    interface{}     `mapstructure:"environments"`
	DeployDefaults DeployManifest  `yaml:"deployDefaults"`
}

type ConfigDeployers struct {
	DeployerGgn ConfigDeployerGgn `mapstructure:"ggn"`
}

type ConfigDeployerGgn struct {
	GitBuildtoolsUrl string `yaml:"gitBuildtoolsUrl"`
	WorkPath         string `yaml:"workPath,omitempty"`
	User             string `yaml:"user,omitempty"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := new(Config)

	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".contactkey"))
	viper.AddConfigPath(configPath)

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
