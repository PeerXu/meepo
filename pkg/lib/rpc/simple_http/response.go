package rpc_simple_http

import "github.com/gin-gonic/gin"

type SimpleDoResponse struct {
	CallResponse []byte `json:"callResponse,omitempty"`
	Error        string `json:"error,omitempty"`
}

func (s *SimpleHttpServer) WriteResponse(c *gin.Context, sc int, res *SimpleDoResponse) {
	c.JSON(sc, res)
}

func (s *SimpleHttpServer) WriteResponseWithError(c *gin.Context, sc int, err error) {
	s.WriteResponse(c, sc, &SimpleDoResponse{
		Error: err.Error(),
	})
}
