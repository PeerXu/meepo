package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
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

	fs := cmd.Flags()
	host, _ := fs.GetString("host")

	targetUrl, err := url.Parse(host)
	if err != nil {
		return err
	}

	targetUrl.Path = "/v1/actions/whoami"

	client := resty.New()
	res, err := client.R().
		Post(targetUrl.String())
	if err != nil {
		return err
	}

	var wr http_api.WhoamiResponse
	if err = json.NewDecoder(bytes.NewReader(res.Body())).Decode(&wr); err != nil {
		return err
	}

	fmt.Println(wr.ID)

	return nil
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
