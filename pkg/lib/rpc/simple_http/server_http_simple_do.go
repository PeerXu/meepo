package rpc_simple_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *SimpleHttpServer) HttpSimpleDo(c *gin.Context) {
	var in SimpleDoRequest
	var err error

	logger := s.GetLogger().WithField("#method", "HttpSimpleDo")

	if err = c.BindJSON(&in); err != nil {
		logger.WithError(err).Errorf("failed to unmarshal request")
		s.WriteResponseWithError(c, http.StatusBadRequest, err)
	}

	logger = logger.WithField("method", in.Method)

	out, err := s.handler.Do(s.context(), in.Method, in.CallRequest)
	if err != nil {
		logger.WithError(err).Errorf("failed to handle")
		s.WriteResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	s.WriteResponse(c, http.StatusOK, &SimpleDoResponse{CallResponse: out})

	logger.Debugf("simple handle")
}
