package commands

import (
	"os"

	"fmt"

	"github.com/remyLemeunier/contactkey/utils"
	"github.com/spf13/cobra"
)

func Execute() {
	configFile, err := utils.ReadFile(utils.DefaultHome)
	if err != nil {
		fmt.Printf("Failed to read default file: %q", err)
		os.Exit(1)
	}

	cfg, err := utils.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Failed load config: %q", err)
		os.Exit(1)
	}

	services, err := cfg.DiscoverServices()
	if err != nil {
		fmt.Printf("Failed to find services: %q", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use: "cck",
	}

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(rollbackCmd)

	deployEnvsCmd := addEnvironmentToCommand(deployCmd, cfg.GlobalEnvironments)
	for _, deployEnvCmd := range deployEnvsCmd {
		addServiceNameToCommand(deployEnvCmd, cfg, services, deployCmd.Name(), deployEnvCmd.Name())
	}

	diffEnvsCmd := addEnvironmentToCommand(diffCmd, cfg.GlobalEnvironments)
	for _, diffEnvCmd := range diffEnvsCmd {
		addServiceNameToCommand(diffEnvCmd, cfg, services, diffCmd.Name(), diffEnvCmd.Name())
	}

	listEnvsCmd := addEnvironmentToCommand(listCmd, cfg.GlobalEnvironments)
	for _, listEnvCmd := range listEnvsCmd {
		addServiceNameToCommand(listEnvCmd, cfg, services, listCmd.Name(), listEnvCmd.Name())
	}

	rollbackEnvsCmd := addEnvironmentToCommand(rollbackCmd, cfg.GlobalEnvironments)
	for _, rollbackEnvCmd := range rollbackEnvsCmd {
		addServiceNameToCommand(rollbackEnvCmd, cfg, services, rollbackCmd.Name(), rollbackEnvCmd.Name())
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
