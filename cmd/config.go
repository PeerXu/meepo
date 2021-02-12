package cmd

import "github.com/spf13/cobra"

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Meepo config subcommand",
	}
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().StringP("config", "c", "~/.meepo/config.yaml", "Location of meepo config file")
}
