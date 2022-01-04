package meepo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
	"golang.org/x/sync/errgroup"

	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/teleportation"
	"github.com/PeerXu/meepo/pkg/transport"
	mconn "github.com/PeerXu/meepo/pkg/util/conn"
	msync "github.com/PeerXu/meepo/pkg/util/sync"
)

const SOCKS5_DOMAIN_SUFFIX = ".mpo"

func isAvailableName(fqdn, suffix string) bool {
	return strings.HasSuffix(fqdn, suffix)
}

type socks5NameResolver struct {
	suffix string
}

func (r *socks5NameResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	if !isAvailableName(name, r.suffix) {
		return nil, nil, &net.DNSError{
			Err:        "no such host",
			Name:       name,
			IsNotFound: true,
		}
	}

	return ctx, net.IPv4(127, 0, 0, 1), nil
}

type socks5Addr struct{}

func (*socks5Addr) Network() string {
	return "socks5"
}

func (*socks5Addr) String() string {
	return "0.0.0.0:0"
}

var SOCKS5ADDR = new(socks5Addr)

func newNewSocks5ServerOption() ofn.Option {
	return ofn.NewOption(map[string]interface{}{})
}

type Socks5Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	Wait() error
}

type socks5Server struct {
	opt    ofn.Option
	logger logrus.FieldLogger
	lis    net.Listener
	eg     errgroup.Group

	meepo               *Meepo
	socks5              *socks5.Server
	dialRequestChannels map[string]chan *teleportation.DialRequest

	domainSuffix string

	lisMtx                 msync.Locker
	dialRequestChannelsMtx msync.Locker
}

func NewSocks5Server(opts ...NewSocks5ServerOption) (Socks5Server, error) {
	var ok bool
	var mp *Meepo
	var logger logrus.FieldLogger

	o := newNewMeepoOption()

	for _, opt := range opts {
		opt(o)
	}

	if logger, ok = o.Get("logger").Inter().(logrus.FieldLogger); !ok {
		logger = logrus.New()
	}

	if mp, ok = o.Get("meepo").Inter().(*Meepo); !ok {
		return nil, fmt.Errorf("require meepo")
	}

	ss := &socks5Server{
		opt:                 o,
		logger:              logger,
		meepo:               mp,
		dialRequestChannels: make(map[string]chan *teleportation.DialRequest),

		domainSuffix: SOCKS5_DOMAIN_SUFFIX,

		lisMtx:                 msync.NewLock(),
		dialRequestChannelsMtx: msync.NewLock(),
	}

	return ss, nil
}

func (s *socks5Server) getLogger() logrus.FieldLogger {
	return s.logger.WithFields(logrus.Fields{
		"#instance": "socks5Server",
		"id":        s.meepo.GetID(),
	})
}

func (s *socks5Server) isAvaiableTransportID(r *socks5.Request) bool {
	return isAvailableName(r.DestAddr.FQDN, s.domainSuffix)
}

func (s *socks5Server) getTransportID(r *socks5.Request) (string, bool) {
	if !s.isAvaiableTransportID(r) {
		return "", false
	}

	return strings.TrimSuffix(r.DestAddr.FQDN, s.domainSuffix), true
}

func (s *socks5Server) mustGetTransportID(r *socks5.Request) string {
	id, _ := s.getTransportID(r)
	return id
}

func (s *socks5Server) channelName(r *socks5.Request) string {
	return fmt.Sprintf("%s:%d", s.mustGetTransportID(r), r.RawDestAddr.Port)
}

func (s *socks5Server) getOrCreateTransport(r *socks5.Request) (transport.Transport, error) {
	peerID := s.mustGetTransportID(r)
	tp, err := s.meepo.getTransport(peerID)
	if err != nil {
		if !errors.Is(err, ErrTransportNotExist) {
			return nil, err
		}

		done := make(chan struct{})
		var doneOnce sync.Once
		tp, err = s.meepo.NewTransport(peerID)
		if err != nil {
			return nil, err
		}
		fn := func(transport.HandleID) {
			doneOnce.Do(func() { close(done) })
		}
		h1 := tp.OnTransportState(transport.TransportStateConnected, fn)
		defer tp.UnsetOnTransportState(transport.TransportStateConnected, h1)
		h2 := tp.OnTransportState(transport.TransportStateFailed, fn)
		defer tp.UnsetOnTransportState(transport.TransportStateFailed, h2)
		h3 := tp.OnTransportState(transport.TransportStateClosed, fn)
		defer tp.UnsetOnTransportState(transport.TransportStateClosed, h3)

		<-done
	}

	return tp, nil
}

func (s *socks5Server) remoteAddr(r *socks5.Request) net.Addr {
	return &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: r.RawDestAddr.Port,
	}
}

func (s *socks5Server) teleport(r *socks5.Request) (chan *teleportation.DialRequest, error) {
	logger := s.getLogger().WithField("#mehtod", "teleport")

	var ts *teleportation.TeleportationSource
	var err error

	tp, err := s.getOrCreateTransport(r)
	if err != nil {
		logger.WithError(err).Debugf("failed to get or create transport")
		return nil, err
	}

	id := tp.PeerID()
	chName := s.channelName(r)
	tpName := fmt.Sprintf("%s:%s", "socks5", chName)
	local := SOCKS5ADDR
	remote := s.remoteAddr(r)

	logger = logger.WithFields(logrus.Fields{
		"peerID":   id,
		"name":     tpName,
		"lnetwork": "socks5",
		"raddr":    remote.String(),
	})

	req := &NewTeleportationRequest{
		Name:          tpName,
		LocalNetwork:  local.Network(),
		LocalAddress:  local.String(),
		RemoteNetwork: remote.Network(),
		RemoteAddress: remote.String(),
	}

	if r.AuthContext.Method == statute.MethodUserPassAuth {
		if req.HashedSecret, err = hashSecret(r.AuthContext.Payload["password"]); err != nil {
			logger.WithError(err).Debugf("failed to hash secret")
			return nil, err
		}
	}

	in := s.meepo.createRequest(id, METHOD_NEW_TELEPORTATION, req)
	out, err := s.meepo.doRequest(in)
	if err != nil {
		logger.WithError(err).Debugf("failed to do request")
		return nil, err
	}

	if err = out.Err(); err != nil {
		logger.WithError(err).Debugf("failed to new teleportation by peer")
		return nil, err
	}

	dialRequests := make(chan *teleportation.DialRequest)
	var closeOnce sync.Once
	closeFn := func() {
		s.dialRequestChannelsMtx.Lock()
		defer s.dialRequestChannelsMtx.Unlock()

		delete(s.dialRequestChannels, chName)
		close(dialRequests)

	}

	ts, err = teleportation.NewTeleportationSource(
		teleportation.WithLogger(s.meepo.getRawLogger()),
		teleportation.WithName(tpName),
		teleportation.WithSource(local),
		teleportation.WithSink(remote),
		teleportation.WithTransport(tp),
		teleportation.SetDialRequestChannel(dialRequests),
		teleportation.WithDoTeleportFunc(func(label string) error {
			innerLogger := logger.WithField("#method", "doTeleportFunc")

			req := &DoTeleportRequest{
				Name:  tpName,
				Label: label,
			}
			in := s.meepo.createRequest(id, METHOD_DO_TELEORT, req)
			out, err := s.meepo.doRequest(in)
			if err != nil {
				innerLogger.WithError(err).Debugf("failed to do request")
				return err
			}

			if err = out.Err(); err != nil {
				innerLogger.WithError(err).Debugf("failed to do teleport by peer")
				return err
			}

			innerLogger.Tracef("do teleport")

			return nil
		}),
		teleportation.WithOnCloseHandler(func() {
			s.meepo.removeTeleportationSource(ts.Name())
			logger.Tracef("remove teleportation source")

			closeOnce.Do(closeFn)
		}),
		teleportation.WithOnErrorHandler(func(err error) {
			s.meepo.removeTeleportationSource(ts.Name())
			logger.WithError(err).Tracef("remove teleportation source")

			closeOnce.Do(closeFn)
		}),
	)
	if err != nil {
		logger.WithError(err).Debugf("failed to new teleportation source")
		return nil, err
	}

	tp.OnTransportState(transport.TransportStateFailed, func(hid transport.HandleID) {
		ts.Close()
		tp.UnsetOnTransportState(transport.TransportStateFailed, hid)
	})

	s.meepo.addTeleportationSource(ts.Name(), ts)
	logger.Tracef("add teleportation source")

	s.dialRequestChannelsMtx.Lock()
	s.dialRequestChannels[chName] = dialRequests
	s.dialRequestChannelsMtx.Unlock()

	logger.Infof("new teleportation source")

	return dialRequests, nil
}

func (s *socks5Server) getOrCreateConnChannel(r *socks5.Request) (chan *teleportation.DialRequest, error) {
	var err error

	cn := s.channelName(r)

	s.dialRequestChannelsMtx.Lock()
	ch, ok := s.dialRequestChannels[cn]
	s.dialRequestChannelsMtx.Unlock()

	if !ok {
		ch, err = s.teleport(r)
		if err != nil {
			return nil, err
		}
	}

	return ch, nil
}

func (s *socks5Server) unsupportedCommandHandler(ctx context.Context, w io.Writer, r *socks5.Request) error {
	logger := s.getLogger().WithField("#method", "unsupportedCommandHandler")

	if err := socks5.SendReply(w, statute.RepCommandNotSupported, nil); err != nil {
		logger.WithError(err).Errorf("failed to send reply")
		return err
	}

	logger.Warningf("Unsupported command")

	return ErrUnsupportedSocks5Command
}

type credentialStore func(password string) bool

func (fn credentialStore) Valid(user, password, userAddr string) bool {
	return fn(password)
}

func (s *socks5Server) asAuthenticator() socks5.Authenticator {
	return socks5.UserPassAuthenticator{
		// TODO: verify socks password
		Credentials: credentialStore(func(password string) bool { return true }),
	}
}

func (s *socks5Server) Start(ctx context.Context) error {
	var err error

	host := cast.ToString(s.opt.Get("host").Inter())
	port := cast.ToString(s.opt.Get("port").Inter())

	opts := []socks5.Option{
		socks5.WithResolver(&socks5NameResolver{suffix: s.domainSuffix}),
		socks5.WithConnectHandle(s.handleConnect),
		socks5.WithBindHandle(s.unsupportedCommandHandler),
		socks5.WithAssociateHandle(s.unsupportedCommandHandler),
		socks5.WithAuthMethods([]socks5.Authenticator{
			s.asAuthenticator(),
			&socks5.NoAuthAuthenticator{},
		}),
	}

	s.socks5 = socks5.NewServer(opts...)

	s.lisMtx.Lock()
	s.lis, err = net.Listen("tcp", net.JoinHostPort(host, port))
	s.lisMtx.Unlock()
	if err != nil {
		return err
	}

	s.eg.Go(func() error {
		return s.socks5.Serve(s.lis)
	})

	return nil
}

func (s *socks5Server) Stop(ctx context.Context) error {
	s.lisMtx.Lock()
	defer s.lisMtx.Unlock()

	if s.lis != nil {
		if err := s.lis.Close(); err != nil {
			return err
		}
		s.lis = nil
	}

	return nil
}

func (s *socks5Server) Wait() error {
	return s.eg.Wait()
}

func (s *socks5Server) handleConnect(ctx context.Context, w io.Writer, r *socks5.Request) error {
	var err error

	logger := s.getLogger().WithField("#mehtod", "handleConnect")

	if !s.isAvaiableTransportID(r) {
		err = ErrNetworkUnreachable
		logger.WithError(err).WithField("fqdn", r.RawDestAddr.FQDN).Errorf("unsupported domain")
		if er := socks5.SendReply(w, statute.RepNetworkUnreachable, nil); er != nil {
			logger.WithError(er).Debugf("failed to send reply")
			return er
		}
		return err
	}

	ch, err := s.getOrCreateConnChannel(r)
	if err != nil {
		logger.WithError(err).Errorf("failed to get or create connChannel")
		if er := socks5.SendReply(w, statute.RepHostUnreachable, nil); er != nil {
			logger.WithError(er).Debugf("failed to send reply")
			return er
		}
		return err
	}

	if err = socks5.SendReply(w, statute.RepSuccess, r.LocalAddr); err != nil {
		logger.WithError(err).Debugf("failed to send reply")
		return err
	}

	conn := mconn.NewRWConn(r.Reader, w, SOCKS5ADDR, s.remoteAddr(r))
	quit := make(chan struct{})

	ch <- teleportation.NewDialRequestWithQuit(conn, quit)

	<-quit

	return nil
}
