package listenerer_http

func (c *HttpConn) Enable() {
	c.enableOnce.Do(func() { close(c.enable) })
}
