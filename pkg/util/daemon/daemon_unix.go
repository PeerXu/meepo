// +build linux darwin

package daemon

import "github.com/VividCortex/godaemon"

func Daemon() {
	godaemon.MakeDaemon(&godaemon.DaemonAttr{})
}
