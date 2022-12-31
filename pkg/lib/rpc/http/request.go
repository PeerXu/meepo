package rpc_http

import (
	"encoding/json"

	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	crypto_interface "github.com/PeerXu/meepo/pkg/lib/crypto/interface"
)

var (
	EMPTY_REQUEST = &DoRequest{}
)

type DoRequest struct {
	Raw         *crypto_interface.Packet `json:"-"`
	Source      []byte                   `json:"-"`
	Destination []byte                   `json:"-"`
	Method      string                   `json:"method"`
	CallRequest []byte                   `json:"callRequest,omitempty"`
}

func (s *HttpServer) UnmarshalDoRequest(in *crypto_core.Packet) (*DoRequest, error) {
	var req DoRequest
	var err error

	logger := s.GetLogger().WithField("#method", "UnmarshalDoRequest")

	if err = s.signer.Verify(in); err != nil {
		logger.WithError(err).Debugf("failed to verify packet")
		return nil, err
	}

	buf, err := s.cryptor.Decrypt(in)
	if err != nil {
		logger.WithError(err).Debugf("failed to decrypt ciphertext")
		return nil, err
	}

	if err = json.Unmarshal(buf, &req); err != nil {
		logger.WithError(err).Debugf("failed to unmarshal plaintext to DoRequest")
		return nil, err
	}

	req.Raw = in
	req.Source = in.Source
	req.Destination = in.Destination

	logger.Tracef("unmarshal DoRequest")

	return &req, nil
}

func (s *HttpCaller) MarshalDoRequest(req *DoRequest) (*crypto_core.Packet, error) {
	logger := s.GetLogger().WithField("#method", "MarshalDoRequest")

	buf, err := json.Marshal(req)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal DoRequest to plaintext")
		return nil, err
	}

	out, err := s.cryptor.Encrypt(req.Destination, buf)
	if err != nil {
		logger.WithError(err).Debugf("failed to encrypt plaintext to packet")
		return nil, err
	}

	if err = s.signer.Sign(out); err != nil {
		logger.WithError(err).Debugf("failed to sign packet")
		return nil, err
	}

	logger.Tracef("marshal DoRequest")

	return out, nil
}
