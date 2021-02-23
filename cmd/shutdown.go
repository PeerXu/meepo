package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	shutdownCmd = &cobra.Command{
		Use:     "shutdown",
		Short:   "Shutdown Meepo",
		Example: "meepo shutdown",
		Aliases: []string{"deny"},
		RunE:    meepoShutdown,
	}
)

func meepoShutdown(cmd *cobra.Command, args []string) error {
	var err error

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	if err = sdk.Shutdown(); err != nil {
		return err
	}

	fmt.Println("Meepo shutting down")

	return nil
}

func init() {
	rootCmd.AddCommand(shutdownCmd)
}
