package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/sdk"
	"github.com/PeerXu/meepo/pkg/util/version"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE:  meepoVersion,
	}
)

func printClientVersion(v *version.V) {
	fmt.Printf("Meepo Client:\n")
	fmt.Printf("  Version:\t\t%v\n", v.Version)
	fmt.Printf("  GoVersion:\t\t%v\n", v.GoVersion)
	fmt.Printf("  GitHash:\t\t%v\n", v.GitHash)
	fmt.Printf("  Bulit:\t\t%v\n", v.Built)
	fmt.Printf("  Platform:\t\t%v\n", v.Platform)
}

func printServerVersion(v *sdk.Version) {
	fmt.Printf("Meepo Server:\n")
	fmt.Printf("  Version:\t\t%v\n", v.Version)
	fmt.Printf("  GoVersion:\t\t%v\n", v.GoVersion)
	fmt.Printf("  GitHash:\t\t%v\n", v.GitHash)
	fmt.Printf("  Bulit:\t\t%v\n", v.Built)
	fmt.Printf("  Platform:\t\t%v\n", v.Platform)
}

func meepoVersion(cmd *cobra.Command, args []string) error {
	printClientVersion(version.Get())
	fmt.Println()

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	sv, err := sdk.Version()
	if err != nil {
		return err
	}
	printServerVersion(sv)

	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
