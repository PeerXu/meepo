package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	encoding_api "github.com/PeerXu/meepo/pkg/api/encoding"
)

type VersionResponse = encoding_api.Version

func (s *HttpServer) Version(c *gin.Context) {
	v := s.meepo.Version()

	res := &VersionResponse{
		Version:   v.Version,
		GoVersion: v.GoVersion,
		GitHash:   v.GitHash,
		Built:     v.Built,
		Platform:  v.Platform,
	}

	c.JSON(http.StatusOK, res)
}
