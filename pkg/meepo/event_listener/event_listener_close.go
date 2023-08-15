package meepo_event_listener

func (el *eventListener) Close() error {
	close(el.queue)
	return nil
}
