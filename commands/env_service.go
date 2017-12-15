package commands

import (
	"github.com/remyLemeunier/contactkey/utils"
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

func addServiceNameToCommand(serviceTree utils.ServiceTree, cmd *cobra.Command, cfg *utils.Config, commandName string, env string) {
	for serviceName, filePath := range serviceTree.Service {
		serviceCmd2 := &cobra.Command{
			Use:   serviceName,
			Short: "Run command for " + serviceName,
			RunE: func(cmd *cobra.Command, args []string) error {
				cckCommand, err := makeInstance(cfg, commandName, cmd.Name(), env, filePath)
				if err != nil {
					return err
				}
				execute(cckCommand)

				return nil
			},
		}

		cmd.AddCommand(serviceCmd2)
	}

	for name, child := range serviceTree.Child {
		serviceCmd1 := &cobra.Command{
			Use:   name,
			Short: "Run command for " + name,
		}
		cmd.AddCommand(serviceCmd1)
		addServiceNameToCommand(child, serviceCmd1, cfg, commandName, env)
	}
}
