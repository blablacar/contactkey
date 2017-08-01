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
