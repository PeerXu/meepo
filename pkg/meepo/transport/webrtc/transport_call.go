package transport_webrtc

import (
	"context"
	"fmt"
	"time"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) defaultCallOptions() option.Option {
	return option.NewOption(map[string]any{
		well_known_option.OPTION_TIMEOUT: 43 * time.Second,
		well_known_option.OPTION_SCOPE:   "usr",
	})
}

func (t *WebrtcTransport) Call(ctx context.Context, method string, req meepo_interface.CallRequest, res meepo_interface.CallResponse, opts ...meepo_interface.CallOption) error {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "Call",
		"method":  method,
	})

	o := option.ApplyWithDefault(t.defaultCallOptions(), opts...)

	timeout, err := well_known_option.GetTimeout(o)
	if err != nil {
		logger.WithError(err).Debugf("bad options")
		return err
	}

	scope, err := well_known_option.GetScope(o)
	if err != nil {
		logger.WithError(err).Debugf("bad options")
		return err
	}
	logger = logger.WithField("scope", scope)

	st := t.pc.ConnectionState()
	if st != webrtc.PeerConnectionStateConnected {
		err = ErrInvalidConnectionStateFn(st.String())
		logger.WithError(err).Debugf("connection state not in connected")
		return err
	}

	data, err := t.marshaler.Marshal(req)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal call request")
		return err
	}

	in := t.newRequest(scope, method, data)
	logger = logger.WithField("session", in.Session)

	if err = t.DoRequest(ctx, in); err != nil {
		logger.WithError(err).Debugf("failed to do request")
		return err
	}

	outs, err := t.WaitResponse(ctx, in)
	if err != nil {
		logger.WithError(err).Debugf("failed to get outs")
		return err
	}
	defer func() {
		lch, ok := t.polls.LoadAndDelete(in.Session)
		if ok {
			lch.Close()
		}
		logger.Tracef("close outs")
	}()

	var out Message
	select {
	case out = <-outs:
		logger.Tracef("read response buffer from outs")
	case <-time.After(timeout):
		err = ErrCallTimeout
		logger.Debugf("call timeout")
		return err
	}

	if out.Error != "" {
		err = fmt.Errorf(out.Error)
		logger.WithError(err).Debugf("call with error")
		return err
	}

	if err = t.unmarshaler.Unmarshal(out.Data, res); err != nil {
		logger.WithError(err).WithFields(logging.Fields{
			"data": string(out.Data),
		}).Debugf("failed to unmarshal call response")
		return err
	}

	logger.Tracef("call")

	return nil
}
