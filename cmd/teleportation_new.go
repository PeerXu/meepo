package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	msdk "github.com/PeerXu/meepo/pkg/sdk"
)

var (
	teleportationNewCmd = &cobra.Command{
		Use:   "new",
		Short: "New teleportation",
		RunE:  meepoTeleportationNew,
	}
)

func meepoTeleportationNew(cmd *cobra.Command, args []string) error {
	var err error

	fs := cmd.Flags()
	id, _ := fs.GetString("id")
	name, _ := fs.GetString("name")
	localNetwork, _ := fs.GetString("local-network")
	localAddress, _ := fs.GetString("local-address")
	remoteNetwork, _ := fs.GetString("remote-network")
	remoteAddress, _ := fs.GetString("remote-address")

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	opt := &msdk.NewTeleportationOption{}

	remote, err := net.ResolveTCPAddr(remoteNetwork, remoteAddress)
	if err != nil {
		return err
	}

	if name != "" {
		opt.Name = name
	}

	if localAddress != "" {
		local, err := net.ResolveTCPAddr(localNetwork, localAddress)
		if err != nil {
			return err
		}
		opt.Source = local
	}

	_, err = sdk.NewTeleportation(id, remote, opt)
	if err != nil {
		return err
	}

	fmt.Printf("Teleportation creating\n")

	return nil
}

func init() {
	teleportationCmd.AddCommand(teleportationNewCmd)

	teleportationNewCmd.PersistentFlags().String("id", "", "Meepo ID")
	teleportationNewCmd.PersistentFlags().String("name", "", "Teleportation name")
	teleportationNewCmd.PersistentFlags().String("local-network", "", "Local listen network")
	teleportationNewCmd.PersistentFlags().String("local-address", "", "Local listen address, if not set, random port will be listen")
	teleportationNewCmd.PersistentFlags().String("remote-network", "", "Remote server network")
	teleportationNewCmd.PersistentFlags().String("remote-address", "", "Remote server address")
}
