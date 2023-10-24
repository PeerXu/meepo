package listenerer_http

func (c *HttpConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}
