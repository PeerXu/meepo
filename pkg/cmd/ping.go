package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	pingCmd = &cobra.Command{
		Use:   "ping [-n <nonce>] <addr>",
		Short: "Ping transport",
		RunE:  meepoPing,
		Args:  cobra.ExactArgs(1),
	}

	pingOptions struct {
		Nonce    uint32
		Count    int
		Interval time.Duration
	}
)

func meepoPing(cmd *cobra.Command, args []string) error {
	target, err := addr.FromString(args[0])
	if err != nil {
		return err
	}

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	fmt.Printf("ping %s:\n", target.String())
	pingAt := time.Now()
	nonce, err := sdk.Ping(target, pingOptions.Nonce)
	if err != nil {
		return err
	}
	fmt.Printf("[%d] from %s nonce=%v time=%v\n", 0, target.String(), nonce, time.Since(pingAt))

	for i := 1; i != pingOptions.Count; i++ {
		time.Sleep(pingOptions.Interval)
		pingAt := time.Now()
		nonce, err := sdk.Ping(target, pingOptions.Nonce+uint32(i))
		if err != nil {
			return err
		}
		fmt.Printf("[%d] from %s nonce=%v time=%v\n", i, target.String(), nonce, time.Since(pingAt))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(pingCmd)

	fs := pingCmd.Flags()

	fs.Uint32VarP(&pingOptions.Nonce, "nonce", "n", 0, "nonce of ping")
	fs.IntVar(&pingOptions.Count, "count", 0, "stop after sending [count] ping request")
	fs.DurationVar(&pingOptions.Interval, "interval", time.Second, "wait [interval] between sending each ping request")
}
