package tools

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"
)

func NewDummyTestLogger() logrus.FieldLogger {
	return NewTestLogger(ioutil.Discard)
}

func NewTestLogger(out io.Writer) logrus.FieldLogger {
	return &TestLogger{
		Logger: log.New(out, "", 0),
	}
}

// TODO: This is necessary because it was requested to use a 'logrus.FieldLogger' (inside the router)
// instead of a 'logrus.StdLogger', that could be mocked by a 'log.New(ioutil.Discard, "", 0)'
type TestLogger struct {
	*log.Logger
}

func (t *TestLogger) WithField(key string, value interface{}) *logrus.Entry { return &logrus.Entry{} }
func (t *TestLogger) WithFields(fields logrus.Fields) *logrus.Entry         { return &logrus.Entry{} }
func (t *TestLogger) WithError(err error) *logrus.Entry                     { return &logrus.Entry{} }
func (t *TestLogger) Debugf(format string, args ...interface{})             {}
func (t *TestLogger) Infof(format string, args ...interface{})              {}
func (t *TestLogger) Warnf(format string, args ...interface{})              {}
func (t *TestLogger) Warningf(format string, args ...interface{})           {}
func (t *TestLogger) Errorf(format string, args ...interface{})             {}
func (t *TestLogger) Debug(args ...interface{})                             {}
func (t *TestLogger) Info(args ...interface{})                              {}
func (t *TestLogger) Print(args ...interface{})                             {}
func (t *TestLogger) Warn(args ...interface{})                              {}
func (t *TestLogger) Warning(args ...interface{})                           {}
func (t *TestLogger) Error(args ...interface{})                             {}
func (t *TestLogger) Debugln(args ...interface{})                           {}
func (t *TestLogger) Infoln(args ...interface{})                            {}
func (t *TestLogger) Warnln(args ...interface{})                            {}
func (t *TestLogger) Warningln(args ...interface{})                         {}
func (t *TestLogger) Errorln(args ...interface{})                           {}
