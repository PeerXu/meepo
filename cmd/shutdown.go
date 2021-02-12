package cmd

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var (
	shutdownCmd = &cobra.Command{
		Use:     "shutdown",
		Short:   "Shutdown Meepo",
		Example: "meepo shutdown",
		Aliases: []string{"deny"},
		RunE:    meepoShutdown,
	}
)

func meepoShutdown(cmd *cobra.Command, args []string) error {
	var err error

	fs := cmd.Flags()
	host, _ := fs.GetString("host")

	targetUrl, err := url.Parse(host)
	if err != nil {
		return err
	}

	targetUrl.Path = "/v1/actions/shutdown"

	client := resty.New()
	_, err = client.R().
		SetBody(map[string]interface{}{}).
		Post(targetUrl.String())
	if err != nil {
		return err
	}

	fmt.Println("Meepo shutting down")

	return nil
}

func init() {
	rootCmd.AddCommand(shutdownCmd)
}
