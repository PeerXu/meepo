package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	pingCmd = &cobra.Command{
		Use:   "ping <id>",
		Short: "Send ping request to peer meepo",
		RunE:  meepoPing,
		Args:  cobra.ExactArgs(1),
	}
)

func meepoPing(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		return fmt.Errorf("Require id")
	}
	id := args[0]

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	if err = sdk.Ping(id); err != nil {
		return err
	}

	fmt.Println("Pong")

	return nil
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
