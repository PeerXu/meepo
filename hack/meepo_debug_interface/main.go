package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/cpl"
)

type TransportStateChangeRecord struct {
	ID                 uint64    `json:"id"`
	HappenedAt         time.Time `json:"happenedAt"`
	Host               string    `json:"host"`
	Target             string    `json:"target"`
	Session            string    `json:"session"`
	State              string    `json:"state"`
	CommonPrefixLength int       `json:"cpl"`
}

type TransportStateChangeRecorder struct {
	currentID  uint64
	tape       []TransportStateChangeRecord
	tapeBus    chan TransportStateChangeRecord
	playersMtx sync.Mutex
	players    map[string]*TransportStateChangePlayer
}

func NewTransportStateChangeRecorder() (*TransportStateChangeRecorder, error) {
	rr := &TransportStateChangeRecorder{
		tape:    make([]TransportStateChangeRecord, 0),
		tapeBus: make(chan TransportStateChangeRecord),
		players: make(map[string]*TransportStateChangePlayer),
	}
	go rr.mainloop()
	return rr, nil
}

func (rr *TransportStateChangeRecorder) AddRecord(r TransportStateChangeRecord) {
	r.ID = atomic.AddUint64(&rr.currentID, 1)
	rr.tape = append(rr.tape, r)
	rr.tapeBus <- r
	fmt.Println("send to bus")
}

func (rr *TransportStateChangeRecorder) NewPlayer() (*TransportStateChangePlayer, error) {
	id := lo.RandomString(8, lo.AlphanumericCharset)
	pr := &TransportStateChangePlayer{
		id:       id,
		receiver: make(chan TransportStateChangeRecord, 16),
	}

	rr.playersMtx.Lock()
	defer rr.playersMtx.Unlock()
	rr.players[id] = pr

	return pr, nil
}

func (rr *TransportStateChangeRecorder) ClosePlayer(id string) error {
	rr.playersMtx.Lock()
	defer rr.playersMtx.Unlock()
	p, ok := rr.players[id]
	if !ok {
		return fmt.Errorf("not found")
	}

	if err := p.Close(); err != nil {
		return err
	}

	delete(rr.players, id)

	return nil
}

func (rr *TransportStateChangeRecorder) mainloop() {
	for r := range rr.tapeBus {
		fmt.Println("read from bus")
		rr.playersMtx.Lock()
		for _, pr := range rr.players {
			pr.receiver <- r
			fmt.Printf("send to player:%s\n", pr.ID())
		}
		rr.playersMtx.Unlock()
	}
}

type TransportStateChangePlayer struct {
	id       string
	receiver chan TransportStateChangeRecord
}

func (pr *TransportStateChangePlayer) ID() string { return pr.id }

func (pr *TransportStateChangePlayer) Close() error {
	close(pr.receiver)
	return nil
}

func (pr *TransportStateChangePlayer) Play() TransportStateChangeRecord {
	r := <-pr.receiver
	fmt.Printf("read from player:%s\n", pr.ID())
	return r
}

func main() {
	rr, err := NewTransportStateChangeRecorder()
	if err != nil {
		panic(err)
	}
	upgrader := websocket.Upgrader{}
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
		req.CommonPrefixLength = cpl.CommonPrefixLen(host, target)
		fmt.Printf("[%s] %s %s [session=%s, cpl=%d] => %s\n", req.HappenedAt.Format(time.RFC3339Nano), req.Host, req.Target, req.Session, req.CommonPrefixLength, req.State)

		rr.AddRecord(req)
	})
	v1Action.Handle(http.MethodGet, "/monitor_transport_state_change", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, c.Error(err))
			return
		}
		defer conn.Close()
		p, err := rr.NewPlayer()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, c.Error(err))
			return
		}
		defer rr.ClosePlayer(p.ID())

		for {
			r := p.Play()
			if err := conn.WriteJSON(r); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, c.Error(err))
				return
			}
		}
	})

	http.ListenAndServe("0.0.0.0:8080", r)
}
