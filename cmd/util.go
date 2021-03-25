package cmd

import (
	"encoding/base64"
	"io/ioutil"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/sdk"
	http_sdk "github.com/PeerXu/meepo/pkg/sdk/http"
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

func ParseValue(val string) (string, error) {
	if strings.HasPrefix(val, "base64:") {
		buf, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(val, "base64:"))
		if err != nil {
			return "", err
		}
		return string(buf), nil
	}

	if strings.HasPrefix(val, "file://") {
		buf, err := ioutil.ReadFile(strings.TrimPrefix(val, "file://"))
		if err != nil {
			return "", err
		}
		return string(buf), nil
	}

	return val, nil
}
