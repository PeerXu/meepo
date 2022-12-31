package logging

import (
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/internal/option"
)

type (
	Logger = *logrus.Entry
	Fields = logrus.Fields
)

func NewLogger(opts ...NewLoggerOption) (Logger, error) {
	o := option.ApplyWithDefault(DefaultNewLoggerOptions(), opts...)

	lvlStr, err := GetLevel(o)
	if err != nil {
		return nil, err
	}

	lvl, err := logrus.ParseLevel(lvlStr)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(lvl)

	return logrus.NewEntry(logger), nil
}
