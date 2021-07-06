package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "meepo",
		SilenceUsage: true,
	}
)

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("log-level", "info", "Logging level")
	rootCmd.PersistentFlags().StringP("host", "H", "http://127.0.0.1:12345", "Daemon API base url")
}
