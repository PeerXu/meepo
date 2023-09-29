package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	watchTransportCmd = &cobra.Command{
		Use:     "watch",
		Aliases: []string{"w"},
		Short:   "Watch transport",
		RunE:    meepoWatchTransport,
	}
)

func meepoWatchTransport(cmd *cobra.Command, args []string) error {
	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	tvs, errs, cancel, err := sdk.WatchTransports()
	if err != nil {
		return err
	}
	defer cancel()

	for tv := range tvs {
		fmt.Printf("%s %s\n", tv.Addr, tv.State)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	select {
	case <-s:
		fmt.Fprintf(os.Stderr, "Caught Ctrl+C.")
	case err := <-errs:
		return err
	}

	return nil
}

func init() {
	transportCmd.AddCommand(watchTransportCmd)
}
