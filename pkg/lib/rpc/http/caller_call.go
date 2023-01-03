package rpc_http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_interface "github.com/PeerXu/meepo/pkg/lib/crypto/interface"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (x *HttpCaller) Call(ctx context.Context, method string, req rpc_interface.CallRequest, res rpc_interface.CallResponse, opts ...rpc_interface.CallOption) error {
	logger := x.GetLogger().WithFields(logging.Fields{
		"#method": "Call",
		"method":  method,
	})

	o := option.Apply(opts...)

	dest, err := well_known_option.GetDestination(o)
	if err != nil {
		logger.WithError(err).Debugf("failed to get destination")
		return err
	}

	logger = logger.WithField("destination", addr.Must(addr.FromBytesWithoutMagicCode(dest)).String())

	buf, err := x.marshaler.Marshal(req)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal call request")
		return err
	}

	doReq := &DoRequest{
		Destination: dest,
		Method:      method,
		CallRequest: buf,
	}

	in, err := x.MarshalDoRequest(doReq)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal DoRequest to packet")
		return err
	}

	out, err := x.doRequest(ctx, in)
	if err != nil {
		logger.WithError(err).Debugf("failed to doRequest")
		return err
	}

	doRes, err := x.UnmarshalDoResponse(out)
	if err != nil {
		logger.WithError(err).Debugf("failed to unmarshal packet to DoResponse")
		return err
	}

	if doRes.Error != "" {
		err = errors.New(doRes.Error)
		logger.WithError(err).Debugf("failed to doRequest")
		return err
	}

	if err = x.unmarshaler.Unmarshal(doRes.CallResponse, res); err != nil {
		logger.WithError(err).Debugf("failed to unmarshal call response")
		return err
	}

	logger.Tracef("call")

	return nil
}

func (x *HttpCaller) doRequest(ctx context.Context, in *crypto_interface.Packet) (*crypto_interface.Packet, error) {
	var out crypto_interface.Packet

	urlStr := x.JoinPath("/v1/actions/do")
	buf, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := x.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(buf, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (x *HttpCaller) JoinPath(p string) string {
	u, _ := url.Parse(x.baseURL)
	u.Path = path.Join(u.Path, p)
	return u.String()
}
