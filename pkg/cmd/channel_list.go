package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	listChannelCmd = &cobra.Command{
		Use:     "list <addr>",
		Aliases: []string{"l", "ls"},
		Short:   "List channels by target",
		RunE:    meepoListChannels,
		Args:    cobra.ExactArgs(1),
	}
)

func meepoListChannels(cmd *cobra.Command, args []string) error {
	targetStr := args[0]
	target, err := addr.FromString(targetStr)
	if err != nil {
		return err
	}

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	cvs, err := sdk.ListChannelsByTarget(target)
	if err != nil {
		return err
	}

	sort.Slice(cvs, func(i, j int) bool {
		if cvs[i].Addr != cvs[j].Addr {
			return cvs[i].Addr < cvs[j].Addr
		}

		return cvs[i].ID < cvs[j].ID
	})

	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{"Addr", "ID", "State", "Mode", "Source", "Sink", "Network", "Address"})
	for _, cv := range cvs {
		tb.Append([]string{cv.Addr, fmt.Sprintf("%d", cv.ID), cv.State, cv.Mode, renderBool(cv.IsSource), renderBool(cv.IsSink), cv.SinkNetwork, cv.SinkAddress})
	}

	tb.Render()

	return nil
}

func init() {
	channelCmd.AddCommand(listChannelCmd)
}
