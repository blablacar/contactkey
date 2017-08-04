package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var DefaultHome = filepath.Join(os.Getenv("HOME"), ".contactkey")

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
