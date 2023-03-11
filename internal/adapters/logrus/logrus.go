package logrus

import (
	"github.com/sirupsen/logrus"

	"github.com/CallumKerson/loggerrific"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	logrusLogger := &Logger{
		logrus.New(),
	}
	logrusLogger.SetLevel(logrus.InfoLevel)
	return logrusLogger
}

func (l *Logger) WithField(key string, value interface{}) loggerrific.Entry {
	return l.Logger.WithField(key, value)
}

func (l *Logger) WithFields(fields map[string]interface{}) loggerrific.Entry {
	return l.Logger.WithFields(fields)
}

func (l *Logger) WithError(err error) loggerrific.Entry {
	return l.Logger.WithError(err)
}

func (l *Logger) SetLevelDebug() {
	l.SetLevel(logrus.DebugLevel)
}

func (l *Logger) SetLevelInfo() {
	l.SetLevel(logrus.InfoLevel)
}

func (l *Logger) SetLevelWarn() {
	l.SetLevel(logrus.WarnLevel)
}

func (l *Logger) SetLevelError() {
	l.SetLevel(logrus.ErrorLevel)
}

func (l *Logger) IsDebugEnabled() bool {
	return l.IsLevelEnabled(logrus.DebugLevel)
}

func (l *Logger) IsInfoEnabled() bool {
	return l.IsLevelEnabled(logrus.InfoLevel)
}

func (l *Logger) IsWarnEnabled() bool {
	return l.IsLevelEnabled(logrus.WarnLevel)
}

func (l *Logger) IsErrorEnabled() bool {
	return l.IsLevelEnabled(logrus.ErrorLevel)
}
