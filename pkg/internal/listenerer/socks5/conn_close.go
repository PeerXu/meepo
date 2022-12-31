package listenerer_socks5

func (c *Socks5Conn) Close() error {
	c.closeOnce.Do(func() { c.close <- struct{}{} })
	return nil
}
