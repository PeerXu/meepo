package logging

import (
	"io"
	"os"
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

	fileStr, err := GetFile(o)
	if err != nil {
		return nil, err
	}

	var out io.Writer
	switch fileStr {
	case "stdout", "-":
		out = os.Stdout
	default:
		if out, err = os.OpenFile(fileStr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			return nil, err
		}
	}

	logger := logrus.New()

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = time.RFC3339Nano
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)

	logger.SetLevel(lvl)
	logger.SetOutput(out)

	return logrus.NewEntry(logger), nil
}
