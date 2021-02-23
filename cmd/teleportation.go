package cmd

import "github.com/spf13/cobra"

var (
	teleportationCmd = &cobra.Command{
		Use:   "teleportation",
		Short: "Meepo teleportation subcommand",
	}
)

func init() {
	rootCmd.AddCommand(teleportationCmd)
}
