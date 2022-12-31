package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	closeChannelCmd = &cobra.Command{
		Use:     "close <addr> <id>",
		Aliases: []string{"c"},
		Short:   "Close channel",
		RunE:    meepoCloseChannel,
		Args:    cobra.ExactArgs(2),
	}
)

func meepoCloseChannel(cmd *cobra.Command, args []string) error {
	targetStr := args[0]
	target, err := addr.FromString(targetStr)
	if err != nil {
		return err
	}

	idStr := args[1]
	id, err := strconv.ParseUint(idStr, 10, 16)
	if err != nil {
		return err
	}

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	err = sdk.CloseChannel(target, uint16(id))
	if err != nil {
		return err
	}

	fmt.Println("channel closed")

	return nil
}

func init() {
	channelCmd.AddCommand(closeChannelCmd)
}
