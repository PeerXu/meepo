package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
)

var (
	configSetCmd = &cobra.Command{
		Use:     "set",
		Short:   "Set Meepo config setting",
		Example: "meepo config set <key> <value>",
		RunE:    meepoConfigSet,
	}
)

func meepoConfigSet(cmd *cobra.Command, args []string) error {
	fs := cmd.Flags()

	key, _ := fs.GetString("key")
	value, _ := fs.GetString("value")
	cp, _ := fs.GetString("config")

	if key == "" {
		return fmt.Errorf("require key")
	}

	if value == "" {
		return fmt.Errorf("require value")
	}

	cfg, _, err := config.Load(cp)
	if err != nil {
		return err
	}

	if err = cfg.Set(key, value); err != nil {
		return err
	}

	if err = cfg.Dump(cp); err != nil {
		return err
	}

	return nil
}

func init() {
	configCmd.AddCommand(configSetCmd)

	configSetCmd.PersistentFlags().StringP("key", "k", "", "Config setting key")
	configSetCmd.PersistentFlags().StringP("value", "v", "", "Config setting value")
}
