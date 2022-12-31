//go:build darwin

package config

func darwinConfigPaths() []string {
	return append(commonConfigPaths(),
		"~/.meepo/meepo.yaml",
	)
}

func init() {
	defaultConfigPaths = darwinConfigPaths
}
