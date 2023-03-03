package logging

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/lib/option"
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

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = time.RFC3339Nano
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)

	logger.SetLevel(lvl)

	return logrus.NewEntry(logger), nil
}
