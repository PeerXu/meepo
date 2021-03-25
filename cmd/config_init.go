package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
)

var (
	configInitCmd = &cobra.Command{
		Use:   "init [--overwrite] [<key1>=<value2> ...]",
		Short: "Initial config file",
		RunE:  meepoConfigInit,
	}
)

func meepoConfigInit(cmd *cobra.Command, args []string) error {
	var loaded bool
	var err error

	fs := cmd.Flags()
	cp, _ := fs.GetString("config")
	overwrite, _ := fs.GetBool("overwrite")

	if _, loaded, err = config.Load(cp); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if !overwrite && loaded {
		return fmt.Errorf("Config already existed")
	}

	cfg := config.NewDefaultConfig()

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

	fmt.Println("Meepo config initialized")

	return nil
}

func init() {
	configCmd.AddCommand(configInitCmd)

	configInitCmd.PersistentFlags().Bool("overwrite", false, "Overwrite exists config file")
}
