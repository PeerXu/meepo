package listenerer_socks5

func (c *Socks5Conn) Read(p []byte) (n int, err error) {
	return c.request.Reader.Read(p)
}
