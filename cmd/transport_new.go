package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	transportNewCmd = &cobra.Command{
		Use:     "new <id>",
		Short:   "New transport",
		Aliases: []string{"n"},
		RunE:    meepoTransportNew,
		Args:    cobra.ExactArgs(1),
	}
)

func meepoTransportNew(cmd *cobra.Command, args []string) error {
	var err error

	id := args[0]

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	_, err = sdk.NewTransport(id)
	if err != nil {
		return err
	}

	fmt.Println("Transport creating")

	return nil
}

func init() {
	transportCmd.AddCommand(transportNewCmd)
}
