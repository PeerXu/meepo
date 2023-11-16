package listenerer_http

func (c *HttpConnectConn) Write(p []byte) (int, error) {
	return c.writer.Write(p)
}
