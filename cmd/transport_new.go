package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	transportNewCmd = &cobra.Command{
		Use:   "new",
		Short: "New transport",
		RunE:  meepoTransportNew,
	}
)

func meepoTransportNew(cmd *cobra.Command, args []string) error {
	var err error

	fs := cmd.Flags()
	id, _ := fs.GetString("id")

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

	transportNewCmd.PersistentFlags().String("id", "", "Meepo ID")
}
