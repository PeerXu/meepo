package listenerer_http

func (c *HttpConnectConn) Close() error {
	c.closeOnce.Do(func() { close(c.close) })
	return nil
}

func (c *HttpGetConn) Close() error {
	c.closeOnce.Do(func() {
		c.rd1.Close()
		c.rd2.Close()
		c.wr1.Close()
		c.wr2.Close()
		close(c.close)
	})
	return nil
}

func (c *HttpGetPipeConn) Close() error {
	return c.HttpGetConn.Close()
}
