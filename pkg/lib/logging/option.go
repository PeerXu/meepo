package logging

import "github.com/PeerXu/meepo/pkg/lib/option"

const (
	OPTION_LEVEL = "level"
)

func DefaultNewLoggerOptions() option.Option {
	return option.NewOption(map[string]any{
		OPTION_LEVEL: "info",
	})
}

type NewLoggerOption = option.ApplyOption

var (
	WithLevel, GetLevel = option.New[string](OPTION_LEVEL)
)
