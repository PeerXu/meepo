package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	msdk "github.com/PeerXu/meepo/pkg/sdk"
)

var (
	teleportationNewCmd = &cobra.Command{
		Use:     "new [-n name] [-l local-address] [-s secret] <id> <remote-address>",
		Short:   "New teleportation",
		Aliases: []string{"n"},
		RunE:    meepoTeleportationNew,
		Args:    cobra.ExactArgs(2),
	}
)

func meepoTeleportationNew(cmd *cobra.Command, args []string) error {
	var err error

	fs := cmd.Flags()
	name, _ := fs.GetString("name")
	localAddress, _ := fs.GetString("local-address")
	secret, _ := fs.GetString("secret")

	id := args[0]
	remoteAddress := args[1]
	localNetwork := "tcp"
	remoteNetwork := "tcp"

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

	if secret != "" {
		opt.Secret = secret
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

	teleportationNewCmd.PersistentFlags().StringP("name", "n", "", "Teleportation name")
	teleportationNewCmd.PersistentFlags().StringP("local-address", "l", "", "Local listen address, if not set, random port will be listen")
	teleportationNewCmd.PersistentFlags().StringP("secret", "s", "", "New teleportation with secret")
}
