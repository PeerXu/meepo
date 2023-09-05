package transport_webrtc

func (c *WebrtcChannel) ready() {
	c.readyOnce.Do(func() { close(c.readyCh) })
}

func (c *WebrtcChannel) readyWithError(err error) {
	c.readyOnce.Do(func() {
		c.readyErr = err
		close(c.readyCh)
	})
}
