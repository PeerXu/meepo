package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/util/version"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE:  meepoVersion,
	}
)

func meepoVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("Version: %v\n", version.Version)
	fmt.Printf("GoVersion: %v\n", version.GoVersion)
	fmt.Printf("GitHash: %v\n", version.GitHash)
	fmt.Printf("Bulit: %v\n", version.Built)
	fmt.Printf("Platform: %v/%v\n", runtime.GOOS, runtime.GOARCH)

	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
