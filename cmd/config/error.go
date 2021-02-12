package config

import "fmt"

func UnsupportedTransportNameError(name string) error {
	return fmt.Errorf("Unsupported transport name: %s", name)

}

func UnsupportedSignalingNameError(name string) error {
	return fmt.Errorf("Unsupported signaling name: %s", name)
}

func UnsupportedApiNameError(name string) error {
	return fmt.Errorf("Unsupported api name: %s", name)
}

func UnsupportedGetConfigKeyError(key string) error {
	return fmt.Errorf("Unsupported get config key: %s", key)
}

func UnsupportedSetConfigKeyError(key string) error {
	return fmt.Errorf("Unsupported set config key: %s", key)
}
