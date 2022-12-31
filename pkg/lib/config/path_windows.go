//go:build windows

package config

func windowsConfigPaths() []string {
	return append(commonConfigPaths())
}

func init() {
	defaultConfigPaths = windowsConfigPaths
}
