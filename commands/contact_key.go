package commands

import (
	"fmt"

	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var branch = ""

func Execute() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed load config: %q", err))
		return
	}
	tree, err := cfg.DiscoverServices()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to find services: %q", err))
		return
	}

	rootCmd := &cobra.Command{
		Use: "cck",
	}
	rootCmd.PersistentFlags().StringVarP(&cfg.LogLevel, "loglevel", "L", "warn", "log level")

	verbsCommands := []*cobra.Command{deployCmd, diffCmd, listCmd, rollbackCmd}
	for _, command := range verbsCommands {
		rootCmd.AddCommand(command)
		envsCmd := addEnvironmentToCommand(command, cfg.GlobalEnvironments)
		for _, envCmd := range envsCmd {
			addServiceNameToCommand(tree, envCmd, cfg, command.Name(), envCmd.Name())
		}
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
