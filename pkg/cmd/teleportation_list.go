package cmd

import (
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	listTeleportationsCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List teleportations",
		RunE:    meepoListTeleportations,
		Args:    cobra.NoArgs,
	}
)

func meepoListTeleportations(cmd *cobra.Command, args []string) error {
	var err error

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	tpvs, err := sdk.ListTeleportations()
	if err != nil {
		return err
	}

	sort.Slice(tpvs, func(i, j int) bool {
		if tpvs[i].Addr != tpvs[j].Addr {
			return tpvs[i].Addr < tpvs[j].Addr
		}

		return tpvs[i].ID < tpvs[j].ID
	})

	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{"ID", "Addr", "Mode", "Source.Network", "Source.Address", "Sink.Network", "Sink.Address"})
	for _, tpv := range tpvs {
		tb.Append([]string{tpv.ID, tpv.Addr, tpv.Mode, tpv.SourceNetwork, tpv.SourceAddress, tpv.SinkNetwork, tpv.SinkAddress})
	}
	tb.Render()

	return nil
}

func init() {
	teleportationCmd.AddCommand(listTeleportationsCmd)
}
