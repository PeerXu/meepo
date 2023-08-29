package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
)

var (
	diagnosticCmd = &cobra.Command{
		Use:   "diagnostic",
		Short: "Meepo diagnostic program",
		RunE:  meepoDiagnostic,
		Args:  cobra.NoArgs,
	}
)

func meepoDiagnostic(cmd *cobra.Command, args []string) error {
	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	report, err := sdk.Diagnostic()
	if err != nil {
		return err
	}

	buf, err := json.Marshal(report)
	if err != nil {
		return err
	}

	fmt.Println(string(buf))

	return nil
}

func init() {
	rootCmd.AddCommand(diagnosticCmd)
}
