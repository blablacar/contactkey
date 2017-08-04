package utils

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configFile, err := ReadFile("./testdata/config.yaml")
	if err != nil {
		t.Fatalf("ReadFile failed with err %q", err)
	}

	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig failed with err %q", err)
	}
	if cfg.WorkPath != "/tmp" {
		t.Errorf("Unexpected workPath : %q", cfg.WorkPath)
	}

	if cfg.GlobalEnvironments[0] != "preprod" || cfg.GlobalEnvironments[1] != "prod" {
		t.Error("Issue with global envs.")
	}
}

func TestDiscoverServices(t *testing.T) {
	c := Config{WorkPath: "./testdata"}
	services, err := c.DiscoverServices()
	if err != nil {
		t.Fatalf("DiscoverServices failed with err %q", err)
	}

	if len(services) != 2 {
		t.Errorf("Services'lenght should be 2 instead got %d", len(services))
	}

	if services[0] != "config" {
		t.Errorf("Should be 'config' instead got: %q", services[0])
	}

	if services[1] != "manifest" {
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
