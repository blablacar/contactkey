package commands

import (
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use: "cck",
	}

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(rollbackCmd)

	deployEnvsCmd := addEnvironmentToCommand(deployCmd)
	for _, deployEnvCmd := range deployEnvsCmd {
		addServiceNameToCommand(deployEnvCmd, deployCmd.Name(), deployEnvCmd.Name())
	}

	diffEnvsCmd := addEnvironmentToCommand(diffCmd)
	for _, diffEnvCmd := range diffEnvsCmd {
		addServiceNameToCommand(diffEnvCmd, diffCmd.Name(), diffEnvCmd.Name())
	}

	listEnvsCmd := addEnvironmentToCommand(listCmd)
	for _, listEnvCmd := range listEnvsCmd {
		addServiceNameToCommand(listEnvCmd, listCmd.Name(), listEnvCmd.Name())
	}

	rollbackEnvsCmd := addEnvironmentToCommand(rollbackCmd)
	for _, rollbackEnvCmd := range rollbackEnvsCmd {
		addServiceNameToCommand(rollbackEnvCmd, rollbackCmd.Name(), rollbackEnvCmd.Name())
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
