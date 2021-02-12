package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
)

var (
	configGetCmd = &cobra.Command{
		Use:     "get",
		Short:   "Get Meepo config setting",
		Example: "meepo config get <key>",
		RunE:    meepoConfigGet,
	}
)

func meepoConfigGet(cmd *cobra.Command, args []string) error {
	fs := cmd.Flags()

	key, _ := fs.GetString("key")
	cp, _ := fs.GetString("config")

	if key == "" {
		return fmt.Errorf("require key")
	}

	cfg, _, err := config.Load(cp)
	if err != nil {
		return err
	}

	val, err := cfg.Get(key)
	if err != nil {
		return err
	}

	fmt.Println(val)

	return nil
}

func init() {
	configCmd.AddCommand(configGetCmd)

	configGetCmd.PersistentFlags().StringP("key", "k", "", "Config setting key")
}
