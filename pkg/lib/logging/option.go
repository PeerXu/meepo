package logging

import "github.com/PeerXu/meepo/pkg/lib/option"

const (
	OPTION_LEVEL = "level"
	OPTION_FILE  = "file"
)

func DefaultNewLoggerOptions() option.Option {
	return option.NewOption(map[string]any{
		OPTION_LEVEL: "info",
		OPTION_FILE:  "stdout",
	})
}

type NewLoggerOption = option.ApplyOption

var (
	WithLevel, GetLevel = option.New[string](OPTION_LEVEL)
	WithFile, GetFile   = option.New[string](OPTION_FILE)
)
