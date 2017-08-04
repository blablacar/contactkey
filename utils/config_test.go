package utils

import (
	"testing"
)

const testConfigDir = "testdata"

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig(testConfigDir)
	if err != nil {
		t.Fatalf("LoadConfig failed with err %q", err)
	}
	if cfg == nil {
		t.Fatal("cfg is nil")
	}

	if cfg.WorkPath != "/tmp" {
		t.Errorf("Unexpected WorkPath %q", cfg.WorkPath)
	}

	if cfg.GlobalEnvironments[0] != "preprod" || cfg.GlobalEnvironments[1] != "prod" {
		t.Error("Issue with global envs.")
	}

	if cfg.Deployers.DeployerGgn.SupportedEnvironment["prod"] != "prod-pa3" || cfg.Deployers.DeployerGgn.SupportedEnvironment["preprod"] != "preprod" {
		t.Error("Issue with ggn supported env")
	}

	if cfg.DeployDefaults.Deploy.PodName != "pod-php-webapp-admin-tools" {
		t.Errorf("Unexpected DeployDefaults.Deployment.PodName %q", cfg.DeployDefaults.Deploy.PodName)
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
