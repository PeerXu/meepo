package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	teleportationCloseCmd = &cobra.Command{
		Use:   "close",
		Short: "Close teleportation",
		RunE:  meepoTeleportationClose,
	}
)

func meepoTeleportationClose(cmd *cobra.Command, args []string) error {
	var err error

	fs := cmd.Flags()
	name, _ := fs.GetString("name")

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	if err = sdk.CloseTeleportation(name); err != nil {
		return err
	}

	fmt.Printf("Teleportation closing\n")

	return nil
}

func init() {
	teleportationCmd.AddCommand(teleportationCloseCmd)

	teleportationCloseCmd.PersistentFlags().String("name", "", "Teleportation name")
}
