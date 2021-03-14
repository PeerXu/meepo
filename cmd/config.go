package cmd

import (
	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Meepo config subcommand",
	}
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().StringP("config", "c", config.GetDefaultConfigPath(), "Location of meepo config file")
}
