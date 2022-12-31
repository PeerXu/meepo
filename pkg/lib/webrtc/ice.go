package webrtc

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pion/webrtc/v3"
)

func ParseICEServer(x string) (webrtc.ICEServer, error) {
	var y webrtc.ICEServer

	if !strings.Contains(x, "://") {
		x = strings.Replace(x, ":", "://", 1)
	}

	u, err := url.Parse(x)
	if err != nil {
		return webrtc.ICEServer{}, err
	}
	y.URLs = append(y.URLs, fmt.Sprintf("%s:%s", u.Scheme, u.Host))
	if u.User != nil {
		y.CredentialType = webrtc.ICECredentialTypePassword
		y.Username = u.User.Username()
		y.Credential, _ = u.User.Password()
	}

	return y, nil
}

func ParseICEServers(xs []string) (ys []webrtc.ICEServer, err error) {
	for _, x := range xs {
		y, err := ParseICEServer(x)
		if err != nil {
			return nil, err
		}
		ys = append(ys, y)
	}
	return ys, nil
}
