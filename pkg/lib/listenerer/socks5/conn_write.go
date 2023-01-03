package listenerer_socks5

func (c *Socks5Conn) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}
