package cmd

import "github.com/spf13/cobra"

var (
	transportCmd = &cobra.Command{
		Use:     "transport",
		Aliases: []string{"t"},
		Short:   "Meepo transport subcommand",
	}
)

func init() {
	rootCmd.AddCommand(transportCmd)
}
