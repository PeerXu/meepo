package sync

type ChannelLocker interface {
	Acquire(id int32) error
	Release(id int32) error
	Get(id int32) (chan interface{}, error)
	GetWithUnlock(id int32) (chan interface{}, func(), error)
}

type channelLocker struct {
	chs map[int32]chan interface{}
	mtx Locker
}

func (t *channelLocker) Acquire(id int32) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	ch := make(chan interface{})
	if _, ok := t.chs[id]; ok {
		defer close(ch)
		return ChannelExistError
	}

	t.chs[id] = ch

	return nil
}

func (t *channelLocker) Release(id int32) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	ch, err := t.getNL(id)
	if err != nil {
		return err
	}
	defer close(ch)

	delete(t.chs, id)

	return nil
}

func (t *channelLocker) Get(id int32) (chan interface{}, error) {
	return t.getNL(id)
}

func (t *channelLocker) GetWithUnlock(id int32) (ch chan interface{}, unlock func(), err error) {
	t.mtx.Lock()

	ch, err = t.getNL(id)
	if err != nil {
		defer t.mtx.Unlock()
		return nil, nil, err
	}

	return ch, t.mtx.Unlock, nil
}

func (t *channelLocker) getNL(id int32) (chan interface{}, error) {
	ch, ok := t.chs[id]
	if !ok {
		return nil, ChannelNotExistError
	}

	return ch, nil
}

func NewChannelLocker() ChannelLocker {
	return &channelLocker{
		chs: make(map[int32]chan interface{}),
		mtx: NewLock(),
	}
}
