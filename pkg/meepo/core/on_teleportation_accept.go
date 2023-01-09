package meepo_core

import (
	mio "github.com/PeerXu/meepo/pkg/lib/io"
	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

func (mp *Meepo) onTeleportationAccept(tp Teleportation, conn listenerer_interface.Conn) {
	defer conn.Close() // nolint:errcheck

	ctx := mp.context()
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":         "onTeleportationAccept",
		"teleportationID": tp.ID(),
		"targetAddr":      tp.Addr(),
		"sourceNetwork":   tp.SourceAddr().Network(),
		"sourceAddress":   tp.SourceAddr().String(),
	})

	t, err := mp.GetTransport(ctx, tp.Addr())
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		return
	}

	var sinkNetwork, sinkAddress string
	switch tp.SourceAddr().Network() {
	case "socks5":
		sinkNetwork = conn.RemoteAddr().Network()
		sinkAddress = conn.RemoteAddr().String()
	default:
		sinkNetwork = tp.SinkAddr().Network()
		sinkAddress = tp.SinkAddr().String()
	}
	logger = logger.WithFields(logging.Fields{
		"sinkNetwork": sinkNetwork,
		"sinkAddress": sinkAddress,
	})

	c, err := t.NewChannel(ctx, sinkNetwork, sinkAddress,
		well_known_option.WithMode(tp.Mode()),
	)
	if err != nil {
		logger.WithError(err).Debugf("failed to new channel")
		return
	}
	defer c.Close(ctx) // nolint:errcheck

	if err = c.WaitReady(); err != nil {
		logger.WithError(err).Debugf("failed to wait channel ready")
		return
	}

	done1 := make(chan struct{})
	go func() {
		defer close(done1)
		n, err := mio.Copy(conn, c.Conn())
		logger.WithError(err).WithFields(logging.Fields{
			"from":  "Channel.Conn",
			"to":    "Listenerer.Conn",
			"bytes": n,
		}).Debugf("copy closed")
	}()

	done2 := make(chan struct{})
	go func() {
		defer close(done2)
		n, err := mio.Copy(c.Conn(), conn)
		logger.WithError(err).WithFields(logging.Fields{
			"from":  "Listenerer.Conn",
			"to":    "Channel.Conn",
			"bytes": n,
		}).Debugf("copy closed")
	}()

	select {
	case <-done1:
	case <-done2:
	}

	logger.Tracef("accept done")
}