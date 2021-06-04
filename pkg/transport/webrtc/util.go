package webrtc_transport

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pion/webrtc/v3"
)

func unmarshalICEServer(x string) webrtc.ICEServer {
	var y webrtc.ICEServer

	if !strings.Contains(x, "://") {
		x = strings.Replace(x, ":", "://", 1)
	}

	u, _ := url.Parse(x)
	y.URLs = append(y.URLs, fmt.Sprintf("%s:%s", u.Scheme, u.Host))
	if u.User != nil {
		y.CredentialType = webrtc.ICECredentialTypePassword
		y.Username = u.User.Username()
		y.Credential, _ = u.User.Password()
	}

	return y
}

func unmarshalICEServers(xs []string) []webrtc.ICEServer {
	var ys []webrtc.ICEServer

	for _, x := range xs {
		ys = append(ys, unmarshalICEServer(x))
	}

	return ys
}
