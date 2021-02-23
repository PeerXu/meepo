package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	transportCloseCmd = &cobra.Command{
		Use:   "close",
		Short: "Close transport",
		RunE:  meepoTransportClose,
	}
)

func meepoTransportClose(cmd *cobra.Command, args []string) error {
	var err error

	fs := cmd.Flags()
	id, _ := fs.GetString("id")

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	if err = sdk.CloseTransport(id); err != nil {
		return err
	}

	fmt.Println("Transport closing")

	return nil
}

func init() {
	transportCmd.AddCommand(transportCloseCmd)

	transportCloseCmd.PersistentFlags().String("id", "", "Meepo ID")
}
