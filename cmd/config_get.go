package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
)

var (
	configGetCmd = &cobra.Command{
		Use:   "get <key>",
		Short: "Get Meepo config setting",
		RunE:  meepoConfigGet,
		Args:  cobra.ExactArgs(1),
	}
)

func meepoConfigGet(cmd *cobra.Command, args []string) error {
	fs := cmd.Flags()
	cp, _ := fs.GetString("config")
	key := args[0]

	cfg, _, err := config.Load(cp)
	if err != nil {
		return err
	}

	val, err := cfg.Get(key)
	if err != nil {
		return err
	}

	fmt.Print(val)

	return nil
}

func init() {
	configCmd.AddCommand(configGetCmd)
}
