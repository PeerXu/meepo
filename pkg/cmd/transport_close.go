package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	closeTransportCmd = &cobra.Command{
		Use:     "close <addr>",
		Aliases: []string{"c"},
		Short:   "Close transport",
		RunE:    meepoCloseTransport,
		Args:    cobra.ExactArgs(1),
	}
)

func meepoCloseTransport(cmd *cobra.Command, args []string) error {
	target, err := addr.FromString(args[0])
	if err != nil {
		return err
	}

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	err = sdk.CloseTransport(target)
	if err != nil {
		return err
	}

	fmt.Println("close transport")

	return nil
}

func init() {
	transportCmd.AddCommand(closeTransportCmd)
}
