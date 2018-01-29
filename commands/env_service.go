package commands

import (
	"github.com/blablacar/contactkey/utils"
	"github.com/spf13/cobra"
)

func addEnvironmentToCommand(cmd *cobra.Command, envs []string) map[int]*cobra.Command {
	envCommands := make(map[int]*cobra.Command)
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

func addServiceNameToCommand(cmd *cobra.Command, cfg *utils.Config, services []string, commandName string, env string) map[int]*cobra.Command {
	serviceNameCommands := make(map[int]*cobra.Command)
	for index, service := range services {
		serviceCmd := &cobra.Command{
			Use:   service,
			Short: "Run command for " + service,
			RunE: func(cmd *cobra.Command, args []string) error {
				cckCommand, err := makeInstance(cfg, commandName, cmd.Name(), env)
				if err != nil {
					return err
				}

				return execute(cckCommand)
			},
		}
		cmd.AddCommand(serviceCmd)
		serviceNameCommands[index] = serviceCmd
	}

	return serviceNameCommands
}
