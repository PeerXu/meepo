package cmd

import (
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	listTransportCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List transports",
		RunE:    meepoListTransports,
		Args:    cobra.NoArgs,
	}
)

func meepoListTransports(cmd *cobra.Command, args []string) error {
	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	tvs, err := sdk.ListTransports()
	if err != nil {
		return err
	}

	sort.Slice(tvs, func(i, j int) bool { return tvs[i].Addr < tvs[j].Addr })

	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{"Addr", "State"})

	for _, tv := range tvs {
		tb.Append([]string{tv.Addr, tv.State})
	}

	tb.Render()

	return nil
}

func init() {
	transportCmd.AddCommand(listTransportCmd)
}
