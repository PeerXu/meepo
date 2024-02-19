package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/cpl"
	"github.com/gin-gonic/gin"
)

type TransportStateChangeRecord struct {
	HappenedAt time.Time `json:"happenedAt"`
	Host       string    `json:"host"`
	Target     string    `json:"target"`
	Session    string    `json:"session"`
	State      string    `json:"state"`
}

func main() {
	r := gin.Default()
	v1Action := r.Group("/v1/actions")
	v1Action.Handle(http.MethodPost, "/transport_state_change", func(c *gin.Context) {
		var req TransportStateChangeRecord
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, c.Error(err))
			return
		}
		host, _ := addr.FromString(req.Host)
		target, _ := addr.FromString(req.Target)
		fmt.Printf("[%s] %s %s[session=%s, cpl=%d] => %s\n", req.HappenedAt.Format(time.RFC3339Nano), req.Host, req.Target, req.Session, cpl.CommonPrefixLen(host, target), req.State)
	})

	http.ListenAndServe("0.0.0.0:8080", r)
}
