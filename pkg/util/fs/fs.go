package fs

import (
	"os"
	"path"
)

func EnsureDirectoryExist(p string) (err error) {
	if err = os.MkdirAll(path.Dir(p), 0755); err != nil {
		return
	}

	return
}
