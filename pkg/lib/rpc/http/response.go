package rpc_http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	crypto_interface "github.com/PeerXu/meepo/pkg/lib/crypto/interface"
)

type DoResponse struct {
	CallResponse []byte `json:"callResponse,omitempty"`
	Error        string `json:"error,omitempty"`
}

func (x *HttpServer) MarshalDoResponse(res *DoResponse, req *DoRequest) (*crypto_interface.Packet, error) {
	logger := x.GetLogger().WithField("#method", "MarshalDoResponse")

	buf, err := json.Marshal(res)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal DoResponse to plaintext")
		return nil, err
	}

	out, err := x.cryptor.Encrypt(req.Source, buf)
	if err != nil {
		logger.WithError(err).Debugf("failed to encrypt plaintext to packet")
		return nil, err
	}

	if err = x.signer.Sign(out); err != nil {
		logger.WithError(err).Debugf("failed to sign packet")
		return nil, err
	}

	logger.Tracef("marshal DoResponse")

	return out, nil
}

func (x *HttpServer) WriteResponse(c *gin.Context, sc int, res *DoResponse, req *DoRequest) {
	pkt, err := x.MarshalDoResponse(res, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(sc, pkt)
}

func (x *HttpServer) WriteResponseWithError(c *gin.Context, sc int, err error, req *DoRequest) {
	x.WriteResponse(c, sc, &DoResponse{Error: err.Error()}, req)
}

func (x *HttpCaller) UnmarshalDoResponse(out *crypto_interface.Packet) (*DoResponse, error) {
	var res DoResponse
	var err error

	logger := x.GetLogger().WithField("#method", "UnmarshalDoResponse")

	if err = x.signer.Verify(out); err != nil {
		logger.WithError(err).Debugf("failed to verify packet")
		return nil, err
	}

	buf, err := x.cryptor.Decrypt(out)
	if err != nil {
		logger.WithError(err).Debugf("failed to decrypt ciphertext")
		return nil, err
	}

	if err = json.Unmarshal(buf, &res); err != nil {
		logger.WithError(err).Debugf("failed to unmarshal plaintext to DoResponse")
		return nil, err
	}

	logger.Tracef("unmarshal DoResponse")

	return &res, nil
}
