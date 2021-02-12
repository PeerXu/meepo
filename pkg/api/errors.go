package api

import "fmt"

func UnsupportedServer(name string) error {
	return fmt.Errorf("Unsupported server: %s", name)
}
