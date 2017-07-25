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

	if cfg.DeployDefaults.Deploy.PodName != "pod-php-webapp-admin-tools" {
		t.Errorf("Unexpected DeployDefaults.Deployment.PodName %q", cfg.DeployDefaults.Deploy.PodName)
	}

}
