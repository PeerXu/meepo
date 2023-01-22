package cmd

import "github.com/spf13/cobra"

var (
	testCmd = &cobra.Command{
		Use:          "test",
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(testCmd)
}
