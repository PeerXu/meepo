package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	transportListCmd = &cobra.Command{
		Use:     "list",
		Short:   "List transports",
		Aliases: []string{"ls", "l"},
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Peer",
		"State",
	})

	for _, tp := range tps {
		table.Append([]string{
			tp.PeerID,
			tp.State,
		})
	}

	table.Render()

	return nil
}

func init() {
	transportCmd.AddCommand(transportListCmd)
}
