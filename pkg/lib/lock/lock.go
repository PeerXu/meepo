package lock

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

var enableLockTracing bool

type Locker interface {
	sync.Locker
	TryLock() bool
}

type tracedLock struct {
	mtx    sync.Mutex
	name   string
	id     string
	logger logging.Logger
	lockAt time.Time
}

func NewLock(opts ...NewLockOption) Locker {
	o := option.Apply(opts...)

	if !enableLockTracing {
		return &sync.Mutex{}
	}

	name, err := well_known_option.GetName(o)
	if err != nil {
		panic(err)
	}

	logger, err := logging.NewLogger(logging.WithLevel("trace"))
	if err != nil {
		panic(err)
	}

	return &tracedLock{
		name:   name,
		id:     fmt.Sprintf("%08x", rand.Int31()),
		logger: logger,
	}
}

func (x *tracedLock) Lock() {
	logger := x.GetLogger().WithField("#method", "Lock")
	logger.Tracef("pre-lock")
	x.mtx.Lock()
	x.lockAt = time.Now()
	logger.WithField("lockAt", x.lockAt).Tracef("lock")
}

func (x *tracedLock) Unlock() {
	logger := x.GetLogger().WithField("#method", "Unlock")
	logger.Tracef("pre-unlock")
	unlockAt := time.Now()
	x.mtx.Unlock()
	logger.WithFields(logging.Fields{
		"unlockAt": unlockAt,
		"since":    unlockAt.Sub(x.lockAt),
	}).Tracef("unlock")
}

func (x *tracedLock) TryLock() (ok bool) {
	logger := x.GetLogger().WithField("#method", "TryLock")
	logger.Tracef("pre-try-lock")
	ok = x.mtx.TryLock()
	if !ok {
		logger.Tracef("try-lock fail")
		return
	}
	x.lockAt = time.Now()
	logger.WithField("lockAt", x.lockAt).Tracef("try-lock ok")
	return
}

func init() {
	lockTracingStr := os.Getenv("MPO_EXPERIMENTAL_LOCK_TRACING")
	enableLockTracing, _ = strconv.ParseBool(lockTracingStr)
}
