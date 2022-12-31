package transport_webrtc

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
)

type NewChannelRequest struct {
	Label     string
	ChannelID uint16
	Network   string
	Address   string
	Mode      string
}

type NewChannelResponse struct{}

func (t *WebrtcTransport) newRemoteChannel(ctx context.Context, mode string, label string, channelID uint16, network, address string) error {
	var res NewChannelResponse

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":   "newRemoteChannel",
		"label":     label,
		"channelID": channelID,
		"network":   network,
		"address":   address,
		"mode":      mode,
	})

	if err := t.Call(ctx, "newChannel", &NewChannelRequest{
		Label:     label,
		ChannelID: channelID,
		Network:   network,
		Address:   address,
		Mode:      mode,
	}, &res, well_known_option.WithScope("sys")); err != nil {
		logger.WithError(err).Debugf("failed to new remote channel")
		return err
	}

	logger.Tracef("new remote channel")

	return nil
}

func (t *WebrtcTransport) onNewChannel(ctx context.Context, _req any) (res any, err error) {
	req := _req.(*NewChannelRequest)
	res = &NewChannelResponse{}

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":   "onNewChannel",
		"label":     req.Label,
		"channelID": req.ChannelID,
		"network":   req.Network,
		"address":   req.Address,
	})

	t.tempDataChannelsMtx.Lock()
	defer t.tempDataChannelsMtx.Unlock()

	if err = t.beforeNewChannelHook(t, req.Network, req.Address); err != nil {
		logger.WithError(err).Debugf("failed to before new channel hook")
		return
	}

	if _, found := t.cs[req.ChannelID]; found {
		err = ErrRepeatedChannelID
		logger.WithError(err).Debugf("repeated channel id")
		return
	}

	tdc, found := t.tempDataChannels[req.Label]
	if !found {
		t.tempDataChannels[req.Label] = &tempDataChannel{req: req}
		go t.removeTimeoutTempDataChannel(req.Label)
		logger.Tracef("create temp data channel")
		return
	}

	if tdc.rwc == nil {
		logger.Tracef("wait for data channel open")
		return
	}

	tdc.req = req

	go t.handleNewChannel(req.Label)
	logger.Tracef("on new channel")

	return
}
