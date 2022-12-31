//go:build linux

package config

func linuxConfigPaths() []string {
	return append(commonConfigPaths(),
		"~/.meepo/meepo.yaml",
		"/etc/meepo/meepo.yaml",
	)
}

func init() {
	defaultConfigPaths = linuxConfigPaths
}
