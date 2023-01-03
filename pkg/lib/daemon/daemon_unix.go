//go:build linux || darwin

package daemon

import "github.com/VividCortex/godaemon"

func unixDaemon() {
	godaemon.MakeDaemon(&godaemon.DaemonAttr{}) // nolint:errcheck
}

func init() {
	Daemon = unixDaemon
}
