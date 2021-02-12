package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
)

var (
	configInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Initial config file",
		RunE:  meepoConfigInit,
	}
)

func meepoConfigInit(cmd *cobra.Command, args []string) error {
	var loaded bool
	var err error

	fs := cmd.Flags()
	cp, _ := fs.GetString("config")
	id, _ := fs.GetString("id")
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

	cfg.Meepo.ID = id

	if err = cfg.Dump(cp); err != nil {
		return err
	}

	fmt.Println("Meepo config initialized")

	return nil
}

func init() {
	configCmd.AddCommand(configInitCmd)

	configInitCmd.PersistentFlags().Bool("overwrite", false, "Overwrite config file")
	configInitCmd.PersistentFlags().String("id", "", "Meepo ID")
}
