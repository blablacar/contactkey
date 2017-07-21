package commands

import (
	"github.com/spf13/cobra"
	"os"
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
		addServiceNameToCommand(deployEnvCmd, deployCmd.Name())
	}

	diffEnvsCmd := addEnvironmentToCommand(diffCmd)
	for _, diffEnvCmd := range diffEnvsCmd  {
		addServiceNameToCommand(diffEnvCmd, diffCmd.Name())
	}

	listEnvsCmd := addEnvironmentToCommand(listCmd)
	for _, listEnvCmd := range listEnvsCmd {
		addServiceNameToCommand(listEnvCmd, listCmd.Name())
	}

	rollbackEnvsCmd := addEnvironmentToCommand(rollbackCmd)
	for _, rollbackEnvCmd := range rollbackEnvsCmd {
		addServiceNameToCommand(rollbackEnvCmd, rollbackCmd.Name())
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}