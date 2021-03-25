package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	msdk "github.com/PeerXu/meepo/pkg/sdk"
)

var (
	teleportCmd = &cobra.Command{
		Use:     "teleport [-n <name>] [-l <local-address>] <id> <remote-address>",
		Short:   "New teleportation in easy way",
		Aliases: []string{"tp"},
		RunE:    meepoTeleport,
		Args:    cobra.ExactArgs(2),
	}
)

func meepoTeleport(cmd *cobra.Command, args []string) error {
	fs := cmd.Flags()
	name, _ := fs.GetString("name")
	localAddress, _ := fs.GetString("local-address")
	peerID := args[0]
	remoteAddress := args[1]
	localNetwork := "tcp"
	remoteNetwork := "tcp"

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	var tpOpt msdk.TeleportOption

	remote, err := net.ResolveTCPAddr(remoteNetwork, remoteAddress)
	if err != nil {
		return err
	}

	if localAddress != "" {
		local, err := net.ResolveTCPAddr(localNetwork, localAddress)
		if err != nil {
			return err
		}
		tpOpt.Local = local
	}

	if name != "" {
		tpOpt.Name = name
	}

	local, err := sdk.Teleport(peerID, remote, &tpOpt)
	if err != nil {
		return err
	}

	fmt.Printf("Teleport SUCCESS\n")
	fmt.Printf("Enjoy your teleportation with %s\n", local.String())

	return nil
}

func init() {
	rootCmd.AddCommand(teleportCmd)

	teleportCmd.PersistentFlags().StringP("name", "n", "", "Teleportation name")
	teleportCmd.PersistentFlags().StringP("local-address", "l", "", "Local listen address")
}
