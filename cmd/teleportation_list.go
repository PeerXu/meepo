package cmd

import (
	"fmt"

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

	fmt.Printf("Name\tTransport\tPortal\tSource\tSink\tChannels\n")
	for _, tp := range tps {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%d\n",
			tp.Name,
			tp.Transport.PeerID,
			tp.Portal,
			fmt.Sprintf("%s:%s", tp.Source.Network, tp.Source.Address),
			fmt.Sprintf("%s:%s", tp.Sink.Network, tp.Sink.Address),
			len(tp.DataChannels),
		)
	}

	return nil
}

func init() {
	teleportationCmd.AddCommand(teleportationListCmd)
}
