// +build windows

package daemon

import (
	"fmt"
	"os"
)

// TODO: Support daemon mode on Windows.
func Daemon() {
	fmt.Fprintf(os.Stderr, "Windows not support daemon now, ignore daemon flag\n")
}
