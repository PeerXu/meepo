package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
)

var (
	whoamiCmd = &cobra.Command{
		Use:   "whoami",
		Short: "Get Meepo ID",
		RunE:  meepoWhoami,
		Args:  cobra.NoArgs,
	}
)

func meepoWhoami(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Lookup("identity-file").Changed {
		filename, err := cmd.Flags().GetString("identity-file")
		if err != nil {
			return err
		}
		pubk, _, err := crypto_core.LoadEd25519Key(filename)
		if err != nil {
			return err
		}
		addr, err := addr.FromBytesWithoutMagicCode(pubk)
		if err != nil {
			return err
		}
		fmt.Println(addr.String())
		return nil
	}

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	addr, err := sdk.Whoami()
	if err != nil {
		return err
	}

	fmt.Println(addr.String())

	return nil
}

func init() {
	whoamiCmd.PersistentFlags().StringP("identity-file", "i", "", "Identity file")
	rootCmd.AddCommand(whoamiCmd)
}
