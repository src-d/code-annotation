package service

import (
	"github.com/sirupsen/logrus"
)

// NewLogger returns a logrus Logger
func NewLogger() logrus.FieldLogger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	return logger
}
