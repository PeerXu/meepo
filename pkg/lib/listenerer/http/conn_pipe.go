package listenerer_http

func (c *HttpGetConn) Pipe() (*HttpGetPipeConn, *HttpGetPipeConn) {
	return &HttpGetPipeConn{
			c, c.rd1, c.wr2,
		}, &HttpGetPipeConn{
			c, c.rd2, c.wr1,
		}
}
