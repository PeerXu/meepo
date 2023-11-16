package listenerer_http

import "time"

func (c *HttpGetConn) SetDeadline(t time.Time) error      { return nil }
func (c *HttpGetConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *HttpGetConn) SetWriteDeadline(t time.Time) error { return nil }
