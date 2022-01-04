package http_sdk

import (
	"bytes"
	"encoding/json"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"

	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/sdk"
)

func WithHost(host string) sdk.NewMeepoSDKOption {
	return func(o ofn.Option) {
		o["host"] = host
	}
}

type MeepoSDK struct {
	sdk.BaseMeepoSDK

	opt    ofn.Option
	client *resty.Client

	host string
}

func (t *MeepoSDK) joinPath(p string) (string, error) {
	u, err := url.Parse(t.host)
	if err != nil {
		return "", err
	}

	u.Path = p

	return u.String(), nil
}

func (t *MeepoSDK) doRequest(path string, req interface{}, res interface{}, expectCode int) error {
	u, err := t.joinPath(path)
	if err != nil {
		return err
	}

	out, err := t.client.R().SetBody(req).Post(u)
	if err != nil {
		return err
	}

	if out.StatusCode() != expectCode {
		return sdk.ExtractError(out.Body())
	}

	if res != nil {
		if err = json.NewDecoder(bytes.NewReader(out.Body())).Decode(res); err != nil {
			return err
		}
	}

	return nil
}

func newNewMeepoSDKOption() ofn.Option {
	return ofn.NewOption(map[string]interface{}{
		"host": "http://localhost:12345",
	})
}

func NewMeepoSDK(opts ...sdk.NewMeepoSDKOption) (sdk.MeepoSDK, error) {
	o := newNewMeepoSDKOption()

	for _, opt := range opts {
		opt(o)
	}

	host := cast.ToString(o.Get("host").Inter())

	return &MeepoSDK{
		opt:    o,
		host:   host,
		client: resty.New(),
	}, nil
}

func init() {
	sdk.RegisterNewMeepoSDKFunc("http", NewMeepoSDK)
}
