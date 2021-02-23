package cmd

import (
	"net"

	"github.com/PeerXu/meepo/pkg/sdk"
	http_sdk "github.com/PeerXu/meepo/pkg/sdk/http"
	"github.com/spf13/cobra"
)

func NewHTTPSDK(cmd *cobra.Command) (sdk.MeepoSDK, error) {
	fs := cmd.Flags()
	host, _ := fs.GetString("host")

	sdk, err := sdk.NewMeepoSDK("http", http_sdk.WithHost(host))
	if err != nil {
		return nil, err
	}

	return sdk, nil
}

func MustResolveTCPAddr(network string, address string) net.Addr {
	addr, _ := net.ResolveTCPAddr(network, address)
	return addr
}
