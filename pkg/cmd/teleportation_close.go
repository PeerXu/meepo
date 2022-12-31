package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	closeTeleportationCmd = &cobra.Command{
		Use:     "close <id>",
		Aliases: []string{"c"},
		Short:   "Close teleportation",
		RunE:    meepoCloseTeleportation,
		Args:    cobra.ExactArgs(1),
	}
)

func meepoCloseTeleportation(cmd *cobra.Command, args []string) error {
	var err error

	id := args[0]

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	if err = sdk.CloseTeleportation(id); err != nil {
		return err
	}

	fmt.Println("teleportation closed")

	return nil
}

func init() {
	teleportationCmd.AddCommand(closeTeleportationCmd)
}
