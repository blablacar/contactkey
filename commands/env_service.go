package commands

import (
	"github.com/spf13/cobra"
	"fmt"
)

func addEnvironmentToCommand(cmd *cobra.Command) map[int]*cobra.Command {
	envCommands := make(map[int]*cobra.Command)
	// @todo Change this into something dynamic
	envs := []string{"prod-pa3", "preprod"}

	for index, env := range envs {
		envCmd := &cobra.Command{
			Use:   env,
			Short: "Run command for " + env,
		}
		cmd.AddCommand(envCmd)
		envCommands[index] = envCmd
	}

	return envCommands
}

func addServiceNameToCommand(cmd *cobra.Command, commandName string) map[int]*cobra.Command {
	serviceNameCommands := make(map[int]*cobra.Command)
	// @todo Change this into something dynamic
	services := []string{"webhooks", "pay-subscription-sesterce"}

	for index, service := range services {
		serviceCmd := &cobra.Command{
			Use:   service,
			Short: "Run command for " + service,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Printf("commandName => %s \n", commandName)
				fmt.Printf("service => %s \n", service)
				fmt.Printf("env => %s \n", cmd.Name())
			},
		}
		cmd.AddCommand(serviceCmd)
		serviceNameCommands[index] = serviceCmd
	}

	return serviceNameCommands
}