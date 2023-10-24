package listenerer_http

func (c *HttpConn) Close() error {
	c.closeOnce.Do(func() { c.close <- struct{}{} })
	return nil
}
