//go:build windows

package daemon

func windowsDaemon() {
	panic("no support daemon for windows")
}

func init() {
	Daemon = windowsDaemon
}
