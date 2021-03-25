package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
)

var (
	configSetCmd = &cobra.Command{
		Use:   "set <key1>=<value1> [<key2>=<value2> ...]",
		Short: "Set Meepo config setting",
		RunE:  meepoConfigSet,
	}
)

func meepoConfigSet(cmd *cobra.Command, args []string) error {
	fs := cmd.Flags()

	cp, _ := fs.GetString("config")

	if len(args) == 0 {
		return fmt.Errorf("Require config(key=value)")
	}

	cfg, _, err := config.Load(cp)
	if err != nil {
		return err
	}

	for _, arg := range args {
		ss := strings.SplitN(arg, "=", 2)
		if len(ss) != 2 {
			return fmt.Errorf("Require config(key=value)")
		}

		key, val := ss[0], ss[1]
		if val, err = ParseValue(val); err != nil {
			return err
		}

		if err = cfg.Set(key, val); err != nil {
			return err
		}
	}

	if err = cfg.Dump(cp); err != nil {
		return err
	}

	return nil
}

func init() {
	configCmd.AddCommand(configSetCmd)
}
