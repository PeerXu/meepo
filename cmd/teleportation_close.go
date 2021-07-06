package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	teleportationCloseCmd = &cobra.Command{
		Use:     "close <name>",
		Short:   "Close teleportation",
		Aliases: []string{"c"},
		RunE:    meepoTeleportationClose,
		Args:    cobra.ExactArgs(1),
	}
)

func meepoTeleportationClose(cmd *cobra.Command, args []string) error {
	var err error

	name := args[0]

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	if err = sdk.CloseTeleportation(name); err != nil {
		return err
	}

	fmt.Printf("Teleportation is closing\n")

	return nil
}

func init() {
	teleportationCmd.AddCommand(teleportationCloseCmd)

	teleportationCloseCmd.PersistentFlags().String("name", "", "Teleportation name")
}
