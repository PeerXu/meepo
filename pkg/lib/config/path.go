package config

func commonConfigPaths() []string {
	return []string{
		"meepo.yaml",
	}
}

var defaultConfigPaths func() []string
