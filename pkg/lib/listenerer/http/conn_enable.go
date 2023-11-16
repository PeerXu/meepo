package listenerer_http

func (c *HttpConnectConn) Enable() {
	c.enableOnce.Do(func() { close(c.enable) })
}

func (c *HttpGetConn) Enable() {
	c.enableOnce.Do(func() { close(c.enable) })
}
