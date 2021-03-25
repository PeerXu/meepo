package config

import "fmt"

type UnsupportedError struct {
	Namespace string
	Name      string
}

func (e UnsupportedError) Error() string {
	return fmt.Sprintf("Unsupported namespace: %v, name: %v", e.Namespace, e.Name)
}

type UnsupportedConfigKeyError struct {
	Key string
}

func (e UnsupportedConfigKeyError) Error() string {
	return fmt.Sprintf("Unsupported config key: %v", e.Key)
}
