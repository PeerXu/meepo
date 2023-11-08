package listenerer_interface

import "time"

type Enabler interface {
	Enable()
	WaitEnabled(timeout time.Duration) error
}
