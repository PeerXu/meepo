package redis_signaling

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"
	"golang.org/x/sync/errgroup"

	"github.com/PeerXu/meepo/pkg/signaling"
	mrand "github.com/PeerXu/meepo/pkg/util/random"
	msync "github.com/PeerXu/meepo/pkg/util/sync"
)

var (
	compress bool
)

type Event struct {
	Name       string                `json:"name"`
	Session    int32                 `json:"session"`
	Descriptor *signaling.Descriptor `json:"descriptor,omitempty"`
	Error      string                `json:"error,omitempty"`
}

func NewEvent(name string, session int32, descriptor *signaling.Descriptor, err error) *Event {
	evt := new(Event)
	evt.Name = name
	evt.Session = session
	if descriptor != nil {
		evt.Descriptor = descriptor
	}
	if err != nil {
		evt.Error = err.Error()
	}
	return evt
}

func zip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		return nil, err
	}
	err = gz.Flush()
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func unzip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func encode(obj interface{}) ([]byte, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	if compress {
		if buf, err = zip(buf); err != nil {
			return nil, err
		}
	}

	out := base64.StdEncoding.EncodeToString(buf)

	return []byte(out), nil
}

func decode(in []byte, obj interface{}) error {
	buf, err := base64.StdEncoding.DecodeString(string(in))
	if err != nil {
		return err
	}

	if compress {
		if buf, err = unzip(buf); err != nil {
			return err
		}
	}

	if err := json.Unmarshal(buf, obj); err != nil {
		return err
	}

	return nil
}

func NewWireEvent(descriptor *signaling.Descriptor) *Event {
	return NewWireEventWithSession(mrand.Random.Int31(), descriptor)
}

func NewWireEventWithSession(session int32, descriptor *signaling.Descriptor) *Event {
	return NewEvent("wire", session, descriptor, nil)
}

func NewWiredEventWithSession(session int32, descriptor *signaling.Descriptor) *Event {
	return NewEvent("wired", session, descriptor, nil)
}

func NewErrorEventWithSession(session int32, err error) *Event {
	return NewEvent("error", session, nil, err)
}

type RedisEngine struct {
	opt             objx.Map
	logger          logrus.FieldLogger
	redisOption     *redis.Options
	redisClientMap  sync.Map
	mainloopQuitMap sync.Map
	events          chan *Event
	eventCache      *lru.ARCCache

	redisClient   *redis.Client
	redisPubsub   *redis.PubSub
	wireHandler   signaling.WireHandler
	channelLocker msync.ChannelLocker
}

func (e *RedisEngine) getLogger() logrus.FieldLogger {
	return e.logger.WithFields(logrus.Fields{
		"#instance": "RedisEngine",
		"id":        e.ID(),
	})
}

func (e *RedisEngine) launch() {
	go e.resolveLoop()
	go e.eventLoop()
}

func (e *RedisEngine) eventLoop() {
	for evt := range e.events {
		key := fmt.Sprintf("%s.%d", evt.Name, evt.Session)

		if _, ok := e.eventCache.Get(key); ok {
			continue
		}

		e.eventCache.Add(key, nil)
		e.onEvent(evt)
	}
}

func (e *RedisEngine) resolveLoop() {
	go e.resolveOnce()
	for range time.Tick(cast.ToDuration(e.opt.Get("resolvePeriod").Inter())) {
		go e.resolveOnce()
	}

}

func (e *RedisEngine) resolveOnce() {
	logger := e.getLogger().WithField("#method", "resolveOnce")

	host, _, err := net.SplitHostPort(e.redisOption.Addr)
	if err != nil {
		logger.WithError(err).Debugf("failed to split redis addr to host and port")
		return
	}
	logger = logger.WithField("host", host)

	addrs, err := net.LookupHost(host)
	if err != nil {
		logger.WithError(err).Debugf("failed to lookup host")
		return
	}

	e.redisClientMap.Range(func(key, val interface{}) bool {
		found := false
		expect := key.(string)
		for _, addr := range addrs {
			if expect == addr {
				found = true
			}
		}
		if !found {
			e.stopMainloop(expect)
		}

		return true
	})

	for _, addr := range addrs {
		if _, ok := e.redisClientMap.Load(addr); !ok {
			go e.mainloop(addr)
		}
	}
}

func (e *RedisEngine) getRedisClient(addr string) (*redis.Client, error) {
	val, ok := e.redisClientMap.Load(addr)
	if !ok {
		return nil, fmt.Errorf("Client not exists")
	}

	cli, ok := val.(*redis.Client)
	if !ok {
		return nil, fmt.Errorf("Unexpected redis client")
	}

	return cli, nil
}

func (e *RedisEngine) healthCheckLoop(addr string) {
	var err error

	logger := e.getLogger().WithFields(logrus.Fields{
		"#method": "healthCheckLoop",
		"addr":    addr,
	})

	defer logger.Tracef("health check loop terminated")
	logger.Tracef("health check loop started")

	ticker := time.NewTicker(cast.ToDuration(e.opt.Get("healthCheckPeriod").Inter()))
	defer ticker.Stop()
	for range ticker.C {
		var eg errgroup.Group
		eg.Go(func() error { return e.healthCheckRedisClient(addr) })
		if err = eg.Wait(); err != nil {
			logger.WithError(err).Debugf("failed to health check redis client")
			return
		}
	}
}

func (e *RedisEngine) healthCheckRedisClient(addr string) error {
	logger := e.getLogger().WithFields(logrus.Fields{
		"#method": "healthCheckRedisClient",
		"addr":    addr,
	})

	cli, err := e.getRedisClient(addr)
	if err != nil {
		defer e.stopMainloop(addr)
		logger.WithError(err).Debugf("failed to get redis client")
		return err
	}

	if err = cli.Ping(e.getContext()).Err(); err != nil {
		defer e.stopMainloop(addr)
		logger.WithError(err).Debugf("failed to ping redis server")
		return err
	}

	return nil
}

func (e *RedisEngine) stopMainloop(addr string) {
	logger := e.getLogger().WithFields(logrus.Fields{
		"#method": "stopMainloop",
		"addr":    addr,
	})

	val, ok := e.mainloopQuitMap.Load(addr)
	if !ok {
		logger.Debugf("mainloop not started")
		return
	}

	close(val.(chan struct{}))

	logger.Tracef("mainloop terminating")
}

func (e *RedisEngine) mainloop(addr string) {
	var err error

	logger := e.getLogger().WithFields(logrus.Fields{
		"#method": "mainloop",
		"addr":    addr,
	})

	opt := *e.redisOption
	_, port, _ := net.SplitHostPort(opt.Addr)
	opt.Addr = net.JoinHostPort(addr, port)
	cli := redis.NewClient(&opt)
	defer func() { logger.WithError(cli.Close()).Tracef("redis client closed") }()

	if err = cli.Ping(e.getContext()).Err(); err != nil {
		logger.WithError(err).Debugf("failed to new redis client")
		return
	}
	logger.Tracef("new redis client")

	e.redisClientMap.Store(addr, cli)
	defer e.redisClientMap.Delete(addr)

	quit := make(chan struct{}, 1)
	e.mainloopQuitMap.Store(addr, quit)
	defer e.mainloopQuitMap.Delete(addr)

	go e.healthCheckLoop(addr)

	ps := cli.Subscribe(e.getContext(), e.parseChannelName(e.ID()))
	logger.Tracef("new redis pubsub")
	defer func() { logger.WithError(ps.Close()).Tracef("redis pubsub closed") }()

	evts := ps.Channel()
	defer logger.Tracef("mainloop terminated")

	logger.Tracef("mainloop start")
	for {
		select {
		case <-quit:
			return
		case msg, ok := <-evts:
			if !ok {
				logger.Debugf("redis pubsub channel closed")
				return
			}
			logger.Tracef("receive message from redis")

			var evt Event
			if err := decode([]byte(msg.Payload), &evt); err != nil {
				logger.WithError(err).Debugf("failed to decode event")
				continue
			}

			e.events <- &evt
		}
	}
}

func (e *RedisEngine) waitWiredEvent(session int32) (*Event, error) {
	ch, err := e.channelLocker.Get(session)
	if err != nil {
		return nil, err
	}

	select {
	case in, ok := <-ch:
		if !ok {
			return nil, SessionChannelClosedError
		}

		evt := in.(*Event)
		if evt.Error != "" {
			return nil, fmt.Errorf(evt.Error)
		}

		return evt, nil
	case <-time.After(cast.ToDuration(e.opt.Get("waitWiredEventTimeout").Inter())):
		return nil, signaling.WireTimeoutError
	}
}

func (e *RedisEngine) parseChannelName(id string) string {
	return "meepo.engine." + id
}

func (e *RedisEngine) ID() string {
	return cast.ToString(e.opt.Get("id").Inter())
}

func (e *RedisEngine) getChannelName() string {
	return e.parseChannelName(e.ID())
}

func (e *RedisEngine) getContext() context.Context {
	return context.TODO()
}

func (e *RedisEngine) sendEvent(id string, evt *Event) error {
	logger := e.getLogger().WithFields(logrus.Fields{
		"#method": "sendEvent",
		"name":    evt.Name,
		"session": evt.Session,
	})

	msg, err := encode(evt)
	if err != nil {
		return err
	}

	var errs []error
	ok := false
	e.redisClientMap.Range(func(key, val interface{}) bool {
		addr := key.(string)
		innerLogger := logger.WithField("addr", addr)
		c, err := e.getRedisClient(addr)
		if err != nil {
			errs = append(errs, err)
			innerLogger.WithError(err).Debugf("failed to get redis client")
			return true
		}

		if err = c.Publish(e.getContext(), e.parseChannelName(id), msg).Err(); err != nil {
			errs = append(errs, err)
			innerLogger.WithError(err).Debugf("failed to send message to redis")
			return true
		}

		ok = true
		innerLogger.Tracef("send message to redis")

		return true
	})

	if !ok {
		if len(errs) > 0 {
			return errs[0]
		} else {
			return NotAvailableRedisClientError
		}
	}

	return nil
}

func (e *RedisEngine) Wire(dst, src *signaling.Descriptor) (*signaling.Descriptor, error) {
	logger := e.getLogger().WithField("#method", "RedisEngine.Wire")

	wireEvt := NewWireEvent(src)
	if err := e.channelLocker.Acquire(wireEvt.Session); err != nil {
		logger.WithError(err).Debugf("failed to acquire session channel")
		return nil, err
	}
	defer func() {
		e.channelLocker.Release(wireEvt.Session)
		logger.Tracef("release session channel")
	}()
	logger.Tracef("acquire session channel")

	var wiredEvt *Event
	var wiredErr error
	var wiredWg sync.WaitGroup

	wiredWg.Add(1)
	go func() {
		wiredEvt, wiredErr = e.waitWiredEvent(wireEvt.Session)
		wiredWg.Done()
	}()

	if err := e.sendEvent(dst.ID, wireEvt); err != nil {
		logger.WithError(err).Debugf("failed to send event")
		return nil, err
	}

	wiredWg.Wait()
	if wiredErr != nil {
		return nil, wiredErr
	}
	logger.Tracef("receive event from session channel")

	return wiredEvt.Descriptor, nil
}

func (e *RedisEngine) onWire(in *Event) {
	var out *Event

	logger := e.getLogger().WithFields(logrus.Fields{
		"#method": "RedisEngine.onWire",
		"source":  in.Descriptor.ID,
		"session": in.Session,
	})

	if e.wireHandler == nil {
		logger.Debugf("ignore event")
		return
	}

	logger.Tracef("wire handler start")
	d, err := e.wireHandler(in.Descriptor)
	if err != nil {
		logger.WithError(err).Debugf("wire handler error")
		out = NewErrorEventWithSession(in.Session, err)
	} else {
		out = NewWiredEventWithSession(in.Session, d)
	}
	logger.Tracef("wire handler done")

	if err = e.sendEvent(in.Descriptor.ID, out); err != nil {
		logger.WithError(err).Debugf("failed to send event")
		return
	}
}

func (e *RedisEngine) onWired(evt *Event) {
	logger := e.getLogger().WithField("#method", "RedisEngine.onWired")

	ch, unlock, err := e.channelLocker.GetWithUnlock(evt.Session)
	if err != nil {
		logger.WithError(err).Debugf("failed to get session channel")
		return
	}
	defer unlock()

	ch <- evt
}

func (e *RedisEngine) onError(evt *Event) {
	logger := e.getLogger().WithField("#method", "RedisEngine.onError")

	ch, unlock, err := e.channelLocker.GetWithUnlock(evt.Session)
	if err != nil {
		logger.WithError(err).Debugf("failed to get session channel")
		return
	}
	defer unlock()

	ch <- evt
}

func (e *RedisEngine) onEvent(evt *Event) {
	logger := e.getLogger().WithFields(logrus.Fields{
		"#method": "RedisEngine.onEvent",
		"event":   evt.Name,
	})

	switch evt.Name {
	case "wire":
		e.onWire(evt)
	case "wired":
		e.onWired(evt)
	case "error":
		e.onError(evt)
	default:
		logger.Debugf("unsupported event")
	}
}

func (e *RedisEngine) OnWire(h signaling.WireHandler) {
	e.wireHandler = h
}

func (e *RedisEngine) Close() error {
	e.redisClientMap.Range(func(key, val interface{}) bool {
		addr := key.(string)
		e.stopMainloop(addr)
		return true
	})

	return nil
}

func NewRedisEngine(opts ...signaling.NewEngineOption) (signaling.Engine, error) {
	var ok bool
	var logger logrus.FieldLogger

	o := DefaultEngineOption()

	for _, opt := range opts {
		opt(o)
	}

	ro, err := redis.ParseURL(cast.ToString(o.Get("url").Inter()))
	if err != nil {
		return nil, err
	}

	if logger, ok = o.Get("logger").Inter().(logrus.FieldLogger); !ok {
		logger = logrus.New()
	}

	eventCache, err := lru.NewARC(64)
	if err != nil {
		return nil, err
	}

	e := &RedisEngine{
		opt:           o,
		redisOption:   ro,
		logger:        logger,
		events:        make(chan *Event),
		eventCache:    eventCache,
		channelLocker: msync.NewChannelLocker(),
	}
	go e.launch()

	return e, nil
}

func init() {
	compress = true
	signaling.RegisterNewEngineFunc("redis", NewRedisEngine)
}
