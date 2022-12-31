package config

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

func FindConfigPath(x string) string {
	var err error
	var ps []string
	if x != "" {
		ps = append(ps, x)
	}
	ps = append(ps, defaultConfigPaths()...)
	for _, p := range ps {
		p, err = homedir.Expand(p)
		if err != nil {
			continue
		}

		_, err = os.Stat(p)
		if err == nil {
			return p
		}
	}
	return ""
}
