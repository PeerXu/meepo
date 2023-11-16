package listenerer_http

func (c *HttpConnectConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}
