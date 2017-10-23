package utils

import (
	"github.com/spf13/viper"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	viper.AddConfigPath("./testdata")
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed with err %q", err)
	}
	if cfg.WorkPath != "/tmp" {
		t.Errorf("Unexpected workPath : %q", cfg.WorkPath)
	}

	if cfg.GlobalEnvironments[0] != "preprod" || cfg.GlobalEnvironments[1] != "prod" {
		t.Error("Issue with global envs.")
	}

	if cfg.StashConfig.Url != "url" {
		t.Errorf("Url exptected url got %q", cfg.StashConfig.Sha1MaxSize)
	}

	if cfg.StashConfig.Password != "password" {
		t.Errorf("Password exptected password got %q", cfg.StashConfig.Sha1MaxSize)
	}

	if cfg.StashConfig.User != "user" {
		t.Errorf("User exptected user got %q", cfg.StashConfig.Sha1MaxSize)
	}

	if cfg.StashConfig.Sha1MaxSize != 7 {
		t.Errorf("Sha1MaxSize exptected 7 got %q", cfg.StashConfig.Sha1MaxSize)
	}

	if cfg.NexusConfig.Repository != "repository" {
		t.Errorf("repository exptected got %q", cfg.NexusConfig.Repository)
	}

	if cfg.NexusConfig.Repository != "repository" {
		t.Errorf("repository exptected got %q", cfg.NexusConfig.Repository)
	}

	if cfg.NexusConfig.Url != "127.0.0.1" {
		t.Errorf("127.0.0.1 exptected got %q", cfg.NexusConfig.Url)
	}

	if cfg.NexusConfig.Group != "group" {
		t.Errorf("group exptected got %q", cfg.NexusConfig.Group)
	}

	if cfg.SlackConfig.Url != "127.0.0.1" {
		t.Errorf("127.0.0.1 exptected got %q", cfg.SlackConfig.Url)
	}

	if cfg.SlackConfig.Token != "token" {
		t.Errorf("token exptected got %q", cfg.SlackConfig.Token)
	}

	if cfg.FileLockConfig.FilePath != "/tmp" {
		t.Errorf("/tmp exptected got %q", cfg.FileLockConfig.FilePath)
	}

	if cfg.ScreenMandatory != true {
		t.Error("ScreenMandatory was expected to be true")
	}
}

func TestDiscoverServices(t *testing.T) {
	c := Config{WorkPath: "./testdata"}
	services, err := c.DiscoverServices()
	if err != nil {
		t.Fatalf("DiscoverServices failed with err %q", err)
	}

	if len(services) != 3 {
		t.Errorf("Services'lenght should be 3 instead got %d", len(services))
	}

	if services[0] != "config" {
		t.Errorf("Should be 'config' instead got: %q", services[0])
	}
	if services[1] != "manifest-k8s" {
		t.Errorf("Should be 'manifest' instead got: %q", services[1])
	}

	if services[2] != "manifest" {
		t.Errorf("Should be 'manifest' instead got: %q", services[1])
	}
}

func TestReadFile(t *testing.T) {
	file, err := ReadFile("./testdata/config.yaml")
	if err != nil {
		t.Fatalf("ReadFile failed with err %q", err)
	}

	if file == nil {
		t.Fatal("File shouldn't be empty")
	}
}
