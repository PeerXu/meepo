package sync

import "sync"

type Locker = sync.Locker

func NewLock() Locker {
	return &sync.Mutex{}
}
