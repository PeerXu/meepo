package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	transportCloseCmd = &cobra.Command{
		Use:     "close <name>",
		Short:   "Close transport",
		Aliases: []string{"c"},
		RunE:    meepoTransportClose,
		Args:    cobra.ExactArgs(1),
	}
)

func meepoTransportClose(cmd *cobra.Command, args []string) error {
	var err error

	id := args[0]

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
}
