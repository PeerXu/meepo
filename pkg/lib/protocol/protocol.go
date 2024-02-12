package lib_protocol

import (
	semver "github.com/Masterminds/semver/v3"
	"github.com/samber/lo"
)

const (
	VERSION_STRING         = "v0.1.0"
	UNKNOWN_VERSION_STRING = "v0.0.0-unknown"
)

var (
	VERSION         = lo.Must(semver.NewVersion(VERSION_STRING))
	UNKNOWN_VERSION = lo.Must(semver.NewVersion(UNKNOWN_VERSION_STRING))
)
