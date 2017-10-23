package utils

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	WorkPath           string `mapstructure:"workPath"`
	LogLevel           string
	GlobalEnvironments []string
	ScreenMandatory    bool
	DeployerConfig     `mapstructure:"deployers"`
	VcsConfig          `mapstructure:"sources"`
	Binaries           `mapstructure:"binaries"`
	HookConfig         `mapstructure:"hooks"`
	LockSystemConfig   `mapstructure:"lockSystem"`
	MetricsConfig      `mapstructure:"metrics"`
}

type DeployerConfig struct {
	DeployerGgnConfig `mapstructure:"ggn"`
	DeployerK8sConfig `mapstructure:"k8s"`
}

type K8sEnvironment struct {
	Cluster string `mapstructure:"cluster"`
}

type DeployerK8sConfig struct {
	WorkPath     string                    `mapstructure:"workPath"`
	Environments map[string]K8sEnvironment `mapstructure:"environments"`
	VcsRegexp    string                    `mapstructure:"vcsRegexp"`
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
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Url         string `mapstructure:"url"`
	Sha1MaxSize int    `mapstructure:"sha1MaxSize"`
}

type NewRelicConfig struct {
	Url    string `mapstructure:"url"`
	ApiKey string `mapstructure:"apiKey"`
}

type Binaries struct {
	NexusConfig `mapstructure:"nexus"`
}

type NexusConfig struct {
	Url        string `mapstructure:"url"`
	Repository string `mapstructure:"repository"`
	Group      string `mapstructure:"group"`
}

type HookConfig struct {
	SlackConfig    `mapstructure:"slack"`
	NewRelicConfig `mapstructure:"newRelic"`
}

type SlackConfig struct {
	Url   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

type LockSystemConfig struct {
	FileLockConfig `mapstructure:"fileLock"`
}

type FileLockConfig struct {
	FilePath string `mapstructure:"filePath"`
}

type MetricsConfig struct {
	PrometheusConfig `mapstructure:"prometheus"`
}

type PrometheusConfig struct {
	Url string `mapstructure:"url"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")            // name of config file (without extension)
	viper.AddConfigPath("$HOME/.contactkey") // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.cck")        // call multiple times to add many search paths
	viper.AddConfigPath(".")                 // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		return nil, err
	}
	viper.Unmarshal(&cfg)
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
