package cmd

import "github.com/spf13/cobra"

var (
	teleportationCmd = &cobra.Command{
		Use:          "teleportation",
		Aliases:      []string{"tp"},
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(teleportationCmd)
}
