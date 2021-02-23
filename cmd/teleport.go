package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	msdk "github.com/PeerXu/meepo/pkg/sdk"
)

var (
	teleportCmd = &cobra.Command{
		Use:     "teleport",
		Short:   "New teleportation in easy way",
		Aliases: []string{"tp"},
		RunE:    meepoTeleport,
	}
)

func meepoTeleport(cmd *cobra.Command, args []string) error {
	fs := cmd.Flags()
	peerID, _ := fs.GetString("id")
	name, _ := fs.GetString("name")
	localNetwork, _ := fs.GetString("local-network")
	localAddress, _ := fs.GetString("local-address")
	remoteNetwork, _ := fs.GetString("remote-network")
	remoteAddress, _ := fs.GetString("remote-address")

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

	teleportCmd.PersistentFlags().String("id", "", "Meepo ID")
	teleportCmd.PersistentFlags().String("name", "", "Teleportation name")
	teleportCmd.PersistentFlags().String("local-network", "tcp", "Local listen network")
	teleportCmd.PersistentFlags().String("local-address", "", "Local listen address")
	teleportCmd.PersistentFlags().String("remote-network", "tcp", "Remote server network")
	teleportCmd.PersistentFlags().String("remote-address", "", "Remote server address")
}
