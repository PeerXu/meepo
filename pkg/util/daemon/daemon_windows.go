// +build windows

package daemon

import "fmt"

func Daemon() {
	panic(fmt.Errorf("Unsupported daemon on windows"))
}
