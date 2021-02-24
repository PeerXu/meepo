package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	teleportationListCmd = &cobra.Command{
		Use:     "list",
		Short:   "List teleportations",
		Aliases: []string{"ls"},
		RunE:    meepoTeleportationList,
	}
)

func meepoTeleportationList(cmd *cobra.Command, args []string) error {
	var err error

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	tps, err := sdk.ListTeleportations()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name",
		"Transport",
		"Portal",
		"Source",
		"Sink",
		"Channels",
	})

	for _, tp := range tps {
		table.Append([]string{
			tp.Name,
			tp.Transport.PeerID,
			tp.Portal,
			fmt.Sprintf("%s:%s", tp.Source.Network, tp.Source.Address),
			fmt.Sprintf("%s:%s", tp.Sink.Network, tp.Sink.Address),
			fmt.Sprintf("%d", len(tp.DataChannels)),
		})
	}

	table.Render()

	return nil
}

func init() {
	teleportationCmd.AddCommand(teleportationListCmd)
}
