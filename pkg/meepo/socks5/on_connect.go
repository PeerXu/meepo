package meepo_socks5

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"

	mio "github.com/PeerXu/meepo/pkg/lib/io"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_core "github.com/PeerXu/meepo/pkg/meepo/core"
)

func (ss *Socks5Server) onConnect(ctx context.Context, wr io.Writer, req *socks5.Request) (err error) {
	logger := ss.GetLogger().WithFields(logging.Fields{
		"#method": "onConnect",
	})

	sinkHost, target, err := ss.parseDomain(req.DestAddr.FQDN)
	if err != nil {
		logger.WithError(err).Debugf("failed to parse domain")
		if er := socks5.SendReply(wr, statute.RepNetworkUnreachable, nil); er != nil {
			logger.WithError(er).Debugf("failed to send reply")
		}
		return err
	}
	sinkAddress := fmt.Sprintf("%s:%d", sinkHost, req.DestAddr.Port)

	t, err := ss.mp.GetTransport(ctx, target)
	if err != nil {
		if !errors.Is(err, meepo_core.ErrTransportNotFound) {
			logger.WithError(err).Debugf("failed to get transport")
			if er := socks5.SendReply(wr, statute.RepServerFailure, nil); er != nil {
				logger.WithError(er)
			}
			return err
		}

		t, err = ss.mp.NewTransport(ctx, target)
		if err != nil {
			logger.WithError(err).Debugf("failed to new transport")
			if er := socks5.SendReply(wr, statute.RepServerFailure, nil); er != nil {
				logger.WithError(er)
			}
			return err
		}
	}

	if err = t.WaitReady(); err != nil {
		logger.WithError(err).Debugf("failed to wait transport ready")
		if er := socks5.SendReply(wr, statute.RepServerFailure, nil); er != nil {
			logger.WithError(er)
		}
		return err
	}

	c, err := t.NewChannel(ctx, "tcp", sinkAddress)
	if err != nil {
		logger.WithError(err).Debugf("failed to new channel")
		if er := socks5.SendReply(wr, statute.RepServerFailure, nil); er != nil {
			logger.WithError(er)
		}
		return err
	}
	defer c.Close(ctx)

	if err = c.WaitReady(); err != nil {
		logger.WithError(err).Debugf("failed to wait channel ready")
		if er := socks5.SendReply(wr, statute.RepServerFailure, nil); er != nil {
			logger.WithError(er)
		}
		return err
	}

	if err = socks5.SendReply(wr, statute.RepSuccess, req.LocalAddr); err != nil {
		logger.WithError(err).Debugf("failed to send success reply to socks5 client")
		return err
	}

	done1 := make(chan struct{})
	go func() {
		defer close(done1)
		n, err := mio.Copy(wr, c.Conn())
		logger.WithError(err).WithFields(logging.Fields{
			"from":  "meepo.Channel",
			"to":    "socks5.Writer",
			"bytes": n,
		}).Debugf("copy closed")
	}()

	done2 := make(chan struct{})
	go func() {
		defer close(done2)
		n, err := mio.Copy(c.Conn(), req.Reader)
		logger.WithError(err).WithFields(logging.Fields{
			"from":  "socks5.Request.Reader",
			"to":    "meepo.Channel",
			"bytes": n,
		}).Debugf("copy closed")
	}()

	select {
	case <-done1:
	case <-done2:
	}

	logger.Tracef("connect closed")

	return nil
}
