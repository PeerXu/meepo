package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	transportListCmd = &cobra.Command{
		Use:     "list",
		Short:   "List transports",
		Aliases: []string{"ls"},
		RunE:    meepoTransportList,
	}
)

func meepoTransportList(cmd *cobra.Command, args []string) error {
	var err error

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	tps, err := sdk.ListTransports()
	if err != nil {
		return err
	}

	fmt.Printf("Peer\t\tState\n")
	for _, tp := range tps {
		fmt.Printf("%s\t\t%s\n", tp.PeerID, tp.State)
	}

	return nil
}

func init() {
	transportCmd.AddCommand(transportListCmd)
}
