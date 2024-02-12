package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
	"github.com/PeerXu/meepo/pkg/lib/version"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE:  meepoVersion,
		Args:  cobra.NoArgs,
	}
)

func printClientVersion(v *version.V) {
	fmt.Printf("Meepo Client:\n")
	fmt.Printf("  Version:\t\t%v\n", v.Version)
	fmt.Printf("  GoVersion:\t\t%v\n", v.GoVersion)
	fmt.Printf("  GitHash:\t\t%v\n", v.GitHash)
	fmt.Printf("  Bulit:\t\t%v\n", v.Built)
	fmt.Printf("  Platform:\t\t%v\n", v.Platform)
	fmt.Printf("  Protocol:\t\t%v\n", v.Protocl)
}

func printServerVersion(v *version.V) {
	fmt.Printf("Meepo Server:\n")
	fmt.Printf("  Version:\t\t%v\n", v.Version)
	fmt.Printf("  GoVersion:\t\t%v\n", v.GoVersion)
	fmt.Printf("  GitHash:\t\t%v\n", v.GitHash)
	fmt.Printf("  Bulit:\t\t%v\n", v.Built)
	fmt.Printf("  Platform:\t\t%v\n", v.Platform)
	fmt.Printf("  Protocol:\t\t%v\n", v.Protocl)
}

func meepoVersion(cmd *cobra.Command, args []string) error {
	printClientVersion(version.Get())
	fmt.Println()

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	v, err := sdk.GetVersion()
	if err != nil {
		return err
	}

	printServerVersion(v)

	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
