package http_api

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PeerXu/meepo/pkg/meepo"
)

type TeleportRequest struct {
	ID            string `json:"id" binding:"required"`
	Name          string `json:"name"`
	RemoteNetwork string `json:"remoteNetwork"`
	RemoteAddress string `json:"remoteAddress" binding:"required"`
	LocalNetwork  string `json:"localNetwork"`
	LocalAddress  string `json:"localAddress"`
}

type TeleportResponse struct {
	LocalNetwork string `json:"localNetwork"`
	LocalAddress string `json:"localAddress"`
}

func (s *HttpServer) Teleport(c *gin.Context) {
	var err error
	var req TeleportRequest
	var local net.Addr

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ParseError(err))
		return
	}

	remote, err := net.ResolveTCPAddr(req.RemoteNetwork, req.RemoteAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, ParseError(err))
		return
	}

	var opts []meepo.TeleportOption

	if req.LocalAddress != "" {
		if local, err = net.ResolveTCPAddr(req.LocalNetwork, req.LocalAddress); err != nil {
			c.JSON(http.StatusBadRequest, ParseError(err))
			return
		}

		opts = append(opts, meepo.WithLocalAddress(local))
	}

	if req.Name != "" {
		opts = append(opts, meepo.WithName(req.Name))
	}

	if local, err = s.meepo.Teleport(req.ID, remote, opts...); err != nil {
		c.JSON(http.StatusInternalServerError, ParseError(err))
		return
	}

	res := &TeleportResponse{
		LocalNetwork: local.Network(),
		LocalAddress: local.String(),
	}

	c.JSON(http.StatusOK, res)
}
