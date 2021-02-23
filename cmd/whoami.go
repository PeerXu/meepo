package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	whoamiCmd = &cobra.Command{
		Use:   "whoami",
		Short: "Get Meepo ID",
		RunE:  meepoWhoami,
	}
)

func meepoWhoami(cmd *cobra.Command, args []string) error {
	var err error

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	id, err := sdk.Whoami()
	if err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
