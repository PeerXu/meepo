package http_sdk

import (
	"net"
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
	"github.com/PeerXu/meepo/pkg/sdk"
)

func (t *MeepoSDK) Teleport(peerID string, remote net.Addr, opt *sdk.TeleportOption) (net.Addr, error) {
	req := &http_api.TeleportRequest{
		ID:            peerID,
		RemoteNetwork: remote.Network(),
		RemoteAddress: remote.String(),
	}

	if opt.Local != nil {
		req.LocalNetwork = opt.Local.Network()
		req.LocalAddress = opt.Local.String()
	}

	if opt.Name != "" {
		req.Name = opt.Name
	}

	if opt.Secret != "" {
		req.Secret = opt.Secret
	}

	var res http_api.TeleportResponse

	if err := t.doRequest("/v1/actions/teleport", req, &res, http.StatusOK); err != nil {
		return nil, err
	}

	local, err := net.ResolveTCPAddr(res.LocalNetwork, res.LocalAddress)
	if err != nil {
		return nil, err
	}

	return local, nil
}
