package rpc_http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
)

func (x *HttpServer) HttpDo(c *gin.Context) {
	var in crypto_core.Packet
	var err error

	logger := x.GetLogger().WithField("#method", "HttpDo")

	if err = c.BindJSON(&in); err != nil {
		logger.WithError(err).Debugf("failed to unmarshal request")
		x.WriteResponseWithError(c, http.StatusBadRequest, err, EMPTY_REQUEST)
		return
	}

	logger = logger.WithField("source", addr.Must(addr.FromBytesWithoutMagicCode(in.Source)).String())

	req, err := x.UnmarshalDoRequest(&in)
	if err != nil {
		logger.WithError(err).Errorf("failed to unmarshal do request")
		x.WriteResponseWithError(c, http.StatusBadRequest, err, EMPTY_REQUEST)
		return
	}

	logger = logger.WithField("method", req.Method)

	out, err := x.handler.Do(x.context(), req.Method, req.CallRequest)
	if err != nil {
		logger.WithError(err).Errorf("failed to handle")
		x.WriteResponseWithError(c, http.StatusInternalServerError, err, req)
		return
	}

	x.WriteResponse(c, http.StatusOK, &DoResponse{CallResponse: out}, req)

	logger.Debugf("handle")
}
