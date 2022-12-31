package version

import (
	"fmt"
	"runtime"
)

var (
	Version   string
	GoVersion string
	GitHash   string
	Built     string
)

type V struct {
	Version   string
	GoVersion string
	GitHash   string
	Built     string
	Platform  string
}

func Get() *V {
	return &V{
		Version:   Version,
		GoVersion: GoVersion,
		GitHash:   GitHash,
		Built:     Built,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
