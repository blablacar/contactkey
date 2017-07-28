package commands

import (
	"github.com/spf13/cobra"
)

func addEnvironmentToCommand(cmd *cobra.Command) map[int]*cobra.Command {
	envCommands := make(map[int]*cobra.Command)
	// @todo Change this into something dynamic
	envs := []string{"preprod"}

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

func addServiceNameToCommand(cmd *cobra.Command, commandName string, env string) (map[int]*cobra.Command, error) {
	serviceNameCommands := make(map[int]*cobra.Command)
	// @todo Change this into something dynamic
	services := []string{"airflow"}

	for index, service := range services {
		serviceCmd := &cobra.Command{
			Use:   service,
			Short: "Run command for " + service,
			RunE: func(cmd *cobra.Command, args []string) error {
				cckCommand, err := makeInstance(commandName)
				if err != nil {
					// @todo catch this error
					return err
				}

				fill(cckCommand, cmd.Name(), env)
				execute(cckCommand)

				return nil
			},
		}
		cmd.AddCommand(serviceCmd)
		serviceNameCommands[index] = serviceCmd
	}

	return serviceNameCommands, nil
}
