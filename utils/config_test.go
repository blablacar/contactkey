package utils

import (
	"testing"

	"github.com/spf13/viper"
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

	if len(cfg.PotentialUsername) != 2 {
		t.Fatalf("PotentialUsername length should be 2 instead got: %q", len(cfg.PotentialUsername))
	}

	if cfg.PotentialUsername[0] != "USER" || cfg.PotentialUsername[1] != "bamboo.jira.username" {
		t.Error("Error PotentialUsername doesn't contain the right informations.")
	}
}

func TestDiscoverServices(t *testing.T) {
	c := Config{WorkPath: "./testdata"}
	serviceTree, err := c.DiscoverServices()
	if err != nil {
		t.Fatalf("DiscoverServices failed with err %q", err)
	}

	if len(serviceTree.Service) != 3 {
		t.Errorf("ServiceTree length should be 3 instead got %d", len(serviceTree.Service))
	}

	if _, ok := serviceTree.Service["config"]; !ok {
		t.Error("Key 'config' not found.")
	}

	if _, ok := serviceTree.Service["manifest-k8s"]; !ok {
		t.Error("Key 'manifest-k8s' not found.")
	}

	if _, ok := serviceTree.Service["manifest"]; !ok {
		t.Error("Key 'manifest' not found.")
	}

	if len(serviceTree.Child) != 2 {
		t.Fatalf("ServiceTree child length should be 2 instead got %d", len(serviceTree.Child))
	}

	dirOne, ok := serviceTree.Child["dirOne"]
	if !ok {
		t.Fatal("Key 'dirOne' not found in Child.")
	}

	dirTwo, ok := serviceTree.Child["dirTwo"]
	if !ok {
		t.Fatal("Key 'dirTwo' not found in Child.")
	}

	if len(dirOne.Service) != 2 {
		t.Fatalf("dirOne should have 2 services, found %s", len(dirOne.Service))
	}

	if len(dirOne.Child) != 1 {
		t.Fatalf("dirOne should have 1 child, found %s", len(dirOne.Child))
	}

	if len(dirTwo.Service) != 1 {
		t.Fatalf("dirTwo should have 1 service, found %s", len(dirTwo.Service))
	}

	if len(dirTwo.Child) != 0 {
		t.Fatalf("dirTwo have no child, found %s", len(dirTwo.Child))
	}

	confa, ok := dirOne.Service["confa"]
	if !ok {
		t.Fatal("Key 'confa' not found in dirOne's service")
	}

	if confa != "./testdata/dirOne/confa.yaml" {
		t.Errorf("confa's path should be './testdata/dirOne/confa.yaml' instead found %s", confa)
	}

	confb, ok := dirOne.Service["confb"]
	if !ok {
		t.Fatal("Key 'confb' not found in dirOne's service")
	}

	if confb != "./testdata/dirOne/confb.yaml" {
		t.Errorf("confb's path should be './testdata/dirOne/confa.yaml' instead found %s", confb)
	}

	confd, ok := dirTwo.Service["confd"]
	if !ok {
		t.Fatal("Key 'confd' not found in dirTwo's service")
	}

	if confd != "./testdata/dirTwo/confd.yaml" {
		t.Errorf("confd's path should be './testdata/dirTwo/confd.yaml' instead found %s", confd)
	}

	subDirOne, ok := dirOne.Child["subDirOne"]
	if !ok {
		t.Fatal("subDirOne not found in dirOne child.")
	}

	if len(subDirOne.Child) != 0 {
		t.Fatalf("subDirOne have no child, found %s", len(subDirOne.Child))
	}

	if len(subDirOne.Service) != 1 {
		t.Fatalf("subDirOne should have 1 service, found %s", len(subDirOne.Service))
	}

	confc, ok := subDirOne.Service["confc"]
	if !ok {
		t.Fatal("Key 'confc' not found in subDirOne's service")
	}

	if confc != "./testdata/dirOne/subDirOne/confc.yaml" {
		t.Errorf("confc's path should be './testdata/dirOne/subDirOne/confc.yaml' instead found %s", confc)
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
