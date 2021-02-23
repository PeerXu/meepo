package http_sdk

import (
	"net"
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
	"github.com/PeerXu/meepo/pkg/sdk"
)

func (t *MeepoSDK) NewTeleportation(peerID string, remote net.Addr, opt *sdk.NewTeleportationOption) (*sdk.Teleportation, error) {
	req := &http_api.NewTeleportationRequest{
		PeerID:        peerID,
		RemoteNetwork: remote.Network(),
		RemoteAddress: remote.String(),
	}
	if opt == nil {
		opt = &sdk.NewTeleportationOption{}
	}

	if opt.Name != "" {
		req.Name = opt.Name
	}

	if source := opt.Source; source != nil && source.String() != "" {
		req.LocalNetwork = source.Network()
		req.LocalAddress = source.String()
	}

	var res http_api.NewTeleportationResponse

	if err := t.doRequest("/v1/actions/new_teleportation", req, &res, http.StatusCreated); err != nil {
		return nil, err
	}

	return res.Teleportation, nil
}
