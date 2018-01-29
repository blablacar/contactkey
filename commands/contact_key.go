package commands

import (
	"fmt"

	"github.com/blablacar/contactkey/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var branch = ""

func Execute() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed load config: %q", err))
	}

	services, err := cfg.DiscoverServices()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to find services: %q", err))
	}

	rootCmd := &cobra.Command{
		Use: "cck",
	}

	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "verbose output")
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
		log.Fatalln(err)
	}
}
