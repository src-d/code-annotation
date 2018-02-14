package service

import (
	"github.com/sirupsen/logrus"
)

// NewLogger returns a logrus Logger
func NewLogger(env string) *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	if env == "dev" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.WarnLevel)
	}

	return logger
}
