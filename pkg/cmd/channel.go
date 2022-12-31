package cmd

import "github.com/spf13/cobra"

var (
	channelCmd = &cobra.Command{
		Use:          "channel",
		Aliases:      []string{"c"},
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(channelCmd)
}
