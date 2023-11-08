package listenerer_socks5

func (c *Socks5Conn) Enable() {
	c.enableOnce.Do(func() { close(c.enable) })
}
