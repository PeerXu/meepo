package listenerer_http

func (c *HttpConn) Write(p []byte) (int, error) {
	return c.writer.Write(p)
}
