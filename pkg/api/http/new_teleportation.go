package http_api

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	encoding_api "github.com/PeerXu/meepo/pkg/api/encoding"
	"github.com/PeerXu/meepo/pkg/meepo"
	"github.com/PeerXu/meepo/pkg/teleportation"
)

type NewTeleportationRequest struct {
	PeerID        string `json:"peerID" binding:"required"`
	RemoteNetwork string `json:"remoteNetwork" binding:"required"`
	RemoteAddress string `json:"remoteAddress" binding:"required"`
	Name          string `json:"name,omitempty"`
	LocalNetwork  string `json:"localNetwork,omitempty"`
	LocalAddress  string `json:"localAddress,omitempty"`
	Secret        string `json:"secret,omitempty"`
}

type NewTeleportationResponse struct {
	Teleportation *encoding_api.Teleportation `json:"teleportation"`
}

func (s *HttpServer) NewTeleportation(c *gin.Context) {
	var err error
	var req NewTeleportationRequest
	var tp teleportation.Teleportation

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ParseError(err))
		return
	}

	remote, err := net.ResolveTCPAddr(req.RemoteNetwork, req.RemoteAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, ParseError(err))
		return
	}

	var opts []meepo.NewTeleportationOption
	if req.Name != "" {
		opts = append(opts, meepo.WithName(req.Name))
	}

	if req.LocalAddress != "" {
		local, err := net.ResolveTCPAddr(req.LocalNetwork, req.LocalAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, ParseError(err))
			return
		}

		opts = append(opts, meepo.WithLocalAddress(local))
	}

	if req.Secret != "" {
		opts = append(opts, meepo.WithSecret(req.Secret))
	}

	if tp, err = s.meepo.NewTeleportation(req.PeerID, remote, opts...); err != nil {
		c.JSON(http.StatusInternalServerError, ParseError(err))
		return
	}

	c.JSON(http.StatusCreated, encoding_api.ConvertTeleportation(tp))
}
