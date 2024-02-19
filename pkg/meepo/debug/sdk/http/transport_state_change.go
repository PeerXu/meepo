package meepo_debug_sdk_http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

func (c *Client) TransportStateChange(ctx context.Context, happenedAt, host, target, session, state string) error {
	urlStr, err := url.JoinPath(c.baseUrl, "transport_state_change")
	if err != nil {
		return err
	}

	body := map[string]any{
		"happenedAt": happenedAt,
		"host":       host,
		"target":     target,
		"session":    session,
		"state":      state,
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	rd := bytes.NewReader(buf)
	req, err := http.NewRequest(http.MethodPost, urlStr, rd)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ParseError(res)
	}

	return nil
}
