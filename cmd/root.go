package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "meepo",
		SilenceUsage: true,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().String("log-level", "info", "Logging level")
	rootCmd.PersistentFlags().String("host", "http://127.0.0.1:12345", "Daemon API base url")
}
