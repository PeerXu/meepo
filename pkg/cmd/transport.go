package cmd

import "github.com/spf13/cobra"

var (
	transportCmd = &cobra.Command{
		Use:          "transport",
		Aliases:      []string{"t"},
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(transportCmd)
}
