package rpc_simple_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (c *SimpleHttpCaller) Call(ctx context.Context, method string, req rpc_interface.CallRequest, res rpc_interface.CallResponse, opts ...rpc_interface.CallOption) error {
	logger := c.GetLogger().WithFields(logging.Fields{
		"#method": "Call",
		"method":  method,
	})

	buf, err := c.marshaler.Marshal(req)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal call request")
		return err
	}

	in := &SimpleDoRequest{
		Method:      method,
		CallRequest: buf,
	}

	out, err := c.doRequest(c.context(), in)
	if err != nil {
		logger.WithError(err).Debugf("failed to doRequest")
		return err
	}

	if out.Error != "" {
		err = fmt.Errorf(out.Error)
		logger.WithError(err).Debugf("failed to doRequest")
		return err
	}

	if res != nil {
		if err = c.unmarshaler.Unmarshal(out.CallResponse, res); err != nil {
			logger.WithError(err).Debugf("failed to unmarshal call response")
			return err
		}
	}

	logger.Tracef("call")

	return nil
}

func (c *SimpleHttpCaller) doRequest(ctx context.Context, in *SimpleDoRequest) (*SimpleDoResponse, error) {
	var out SimpleDoResponse
	urlStr := c.JoinPath("/v1/actions/simpleDo")
	buf, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(buf, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *SimpleHttpCaller) JoinPath(p string) string {
	u, _ := url.Parse(c.baseURL)
	u.Path = path.Join(u.Path, p)
	return u.String()
}
