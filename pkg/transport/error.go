package transport

import "fmt"

var (
	DataChannelNotFoundError = fmt.Errorf("DataChannel not found")
)

func UnsupportedTransportError(name string) error {
	return fmt.Errorf("Unsupported transport: %s", name)
}
