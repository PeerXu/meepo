package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
	"github.com/PeerXu/meepo/pkg/lib/dialer"
	listenerer_http "github.com/PeerXu/meepo/pkg/lib/listenerer/http"
	listenerer_socks5 "github.com/PeerXu/meepo/pkg/lib/listenerer/socks5"
)

var (
	newTeleportationCmd = &cobra.Command{
		Use:     "new [--source-network <local-network>] [-l <local-address>] [-m <mode>] <addr> <remote-address>",
		Aliases: []string{"n"},
		Short:   "New teleportation",
		RunE:    meepoNewTeleportation,
		Args:    cobra.ExactArgs(2),
	}

	newTeleportationOptions struct {
		SourceNetwork string
		SourceAddress string
		Mode          string
	}
)

func meepoNewTeleportation(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) < 2 {
		return fmt.Errorf("require addr and remote-address")
	}

	targetStr := args[0]
	target, err := addr.FromString(targetStr)
	if err != nil {
		return err
	}

	sinkAddress := args[1]
	sourceNetwork := newTeleportationOptions.SourceNetwork
	sourceAddress := newTeleportationOptions.SourceAddress
	mode := newTeleportationOptions.Mode

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	switch sourceNetwork {
	case listenerer_http.NAME, listenerer_socks5.NAME:
		sinkAddress = "*"
	}

	tpv, err := sdk.NewTeleportation(target, dialer.NewAddr(sourceNetwork, sourceAddress), dialer.NewAddr("tcp", sinkAddress), mode)
	if err != nil {
		return err
	}

	fmt.Printf("teleportation %s created, listen on %s\n", tpv.ID, tpv.SourceAddress)

	return nil
}

func init() {
	newTeleportationCmd.Flags().StringVar(&newTeleportationOptions.SourceNetwork, "source-network", "tcp", "Source network")
	newTeleportationCmd.Flags().StringVarP(&newTeleportationOptions.SourceAddress, "listen", "l", "127.0.0.1:0", "Listen address")
	newTeleportationCmd.Flags().StringVarP(&newTeleportationOptions.Mode, "mode", "m", "mux", "Teleportation mode [raw, mux, kcp]")

	teleportationCmd.AddCommand(newTeleportationCmd)
}
