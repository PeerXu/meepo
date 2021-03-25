package auth

import (
	"encoding/base64"
	"strings"

	"github.com/stretchr/objx"
)

type OFN = func(objx.Map)

// NewAuthorizeOption

func WithSecret(srtStr string) OFN {
	return func(o objx.Map) {
		var srt []byte
		var err error

		if strings.HasPrefix(srtStr, "base64:") {
			srt, err = base64.StdEncoding.DecodeString(strings.TrimPrefix(srtStr, "base64:"))
			if err != nil {
				panic("bad base64 secret")
			}
		} else {
			srt = []byte(srtStr)
		}

		o["secret"] = srt
	}
}

func WithHashAlgorithm(algo string) OFN {
	return func(o objx.Map) {
		o["hashAlgorithm"] = algo
	}
}

func WithTemplate(tmpl string) OFN {
	return func(o objx.Map) {
		o["template"] = tmpl
	}
}
