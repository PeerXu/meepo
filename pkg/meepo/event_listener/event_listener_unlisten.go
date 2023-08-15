package meepo_event_listener

func (el *eventListener) Unlisten(id string) {
	el.set.Remove(id)
}
