package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Send ping request to peer meepo",
		RunE:  meepoPing,
	}
)

func meepoPing(cmd *cobra.Command, args []string) error {
	var err error

	fs := cmd.Flags()
	id, _ := fs.GetString("id")

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

	pingCmd.PersistentFlags().String("id", "", "Meepo ID")
}
