package utils

import (
	"testing"
)

func TestLoadManifest(t *testing.T) {
	configFile, err := ReadFile("./testdata/manifest.yaml")
	if err != nil {
		t.Fatalf("ReadFile failed with err %q", err)
	}
	manifest, err := LoadManifest(configFile)
	if err != nil {
		t.Fatalf("LoadDeployfile failed with err %q", err)
	}

	if manifest == (&Manifest{}) {
		t.Errorf("Unexpected manifest %q", manifest)
	}

	if manifest.Pod != "pod-webhooks" {
		t.Errorf("Pod should be 'pod-webhooks' instead got  %q", manifest.Pod)
	}

	if manifest.NexusManifest.Artifact != "pod-webhooks" {
		t.Errorf("artifact in the NexusManifest not found got %q", manifest.NexusManifest.Artifact)
	}

	if manifest.SlackManifest.Channel != "channel" {
		t.Errorf("channel in the SlackManifest not found got %q", manifest.SlackManifest.Channel)
	}

	if len(manifest.ExecCommandManifest.OnPreDeploy) != 1 {
		t.Fatalf("OnPredDeploy in the ExecCommandManifest should have a len of 1 instead got %d", len(manifest.ExecCommandManifest.OnPreDeploy))
	}

	if len(manifest.ExecCommandManifest.OnPostDeploy) != 1 {
		t.Fatalf("OnPostDeploy in the ExecCommandManifest should have a len of 1 instead got %d", len(manifest.ExecCommandManifest.OnPostDeploy))
	}

	if manifest.ExecCommandManifest.OnInit[0].Command != "ls" {
		t.Errorf("The OnInit command should be 'ls' instead got %d", manifest.ExecCommandManifest.OnInit[0].Command)
	}

	if len(manifest.ExecCommandManifest.OnInit[0].Args) != 1 {
		t.Fatalf("Args len from OnInit in the ExecCommandManifest should have a len of 1 instead got %d", len(manifest.ExecCommandManifest.OnInit[0].Args))
	}

	if manifest.ExecCommandManifest.OnPreDeploy[0].Command != "ls" {
		t.Errorf("The OnPreDeploy command should be 'ls' instead got %d", manifest.ExecCommandManifest.OnPreDeploy[0].Command)
	}

	if len(manifest.ExecCommandManifest.OnPreDeploy[0].Args) != 1 {
		t.Fatalf("Args len from OnPreDeploy in the ExecCommandManifest should have a len of 1 instead got %d", len(manifest.ExecCommandManifest.OnPreDeploy[0].Args))
	}

	if manifest.ExecCommandManifest.OnPreDeploy[0].Args[0] != "-lah" {
		t.Errorf("Args from OnPreDeploy should have been -lah instead got %d", manifest.ExecCommandManifest.OnPreDeploy[0].Args[0])
	}

	if manifest.ExecCommandManifest.OnPostDeploy[0].Command != "cd /tmp" {
		t.Errorf("The OnPostDeploy command should be 'cd /tmp' instead got %d", manifest.ExecCommandManifest.OnPreDeploy[0].Command)
	}

	if len(manifest.ExecCommandManifest.OnPostDeploy[0].Args) != 0 {
		t.Fatalf("Args len from OnPostDeploy in the ExecCommandManifest should have a len of 0 instead got %d", len(manifest.ExecCommandManifest.OnPostDeploy[0].Args))
	}

	if manifest.ExecCommandManifest.StopOnError != true {
		t.Error("StopOnError for exec manifest expected to be true.")
	}

	if manifest.SlackManifest.StopOnError != false {
		t.Error("StopOnError for slack manifest expected to be false.")
	}
}

func TestDeployerGGNManifest(t *testing.T) {
	d := DeployerGgnManifest{}
	err := d.validate()
	if err == nil {
		t.Error("Validate() on incomplete manifest did not return err")
	}
}

func TestK8sManifest(t *testing.T) {
	configFile, err := ReadFile("./testdata/manifest-k8s.yaml")
	if err != nil {
		t.Fatalf("ReadFile failed with err %q", err)
	}

	manifest, err := LoadManifest(configFile)
	if err != nil {
		t.Fatalf("LoadDeployfile failed with err %q", err)
	}

	if manifest == (&Manifest{}) {
		t.Errorf("Unexpected manifest %q", manifest)
	}

	if manifest.DeployerK8sManifest.Release != "webhooks" {
		t.Errorf("Unexpected manifest %q", manifest)
	}

	if manifest.DeployerK8sManifest.Namespace != "webapps" {
		t.Errorf("Unexpected manifest %q", manifest)
	}
}
