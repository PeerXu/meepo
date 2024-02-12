package version

import (
	"fmt"
	"runtime"

	lib_protocol "github.com/PeerXu/meepo/pkg/lib/protocol"
)

var (
	Version   string
	GoVersion string
	GitHash   string
	Built     string
	Protocol  string
)

type V struct {
	Version   string
	GoVersion string
	GitHash   string
	Built     string
	Platform  string
	Protocl   string
}

func Get() *V {
	return &V{
		Version:   Version,
		GoVersion: GoVersion,
		GitHash:   GitHash,
		Built:     Built,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Protocl:   lib_protocol.VERSION.String(),
	}
}
