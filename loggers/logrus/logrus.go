package loggers

import (
	"github.com/fvdveen/fuzzy"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logrus.FieldLogger
}

func (l logrusLogger) WithField(s string, i interface{}) fuzzy.Logger {
	return logrusLogger{l.FieldLogger.WithField(s, i)}
}

func (l logrusLogger) WithFields(m map[string]interface{}) fuzzy.Logger {
	return logrusLogger{l.FieldLogger.WithFields(m)}
}

// New creates a logrus based Logger
func New(lvl fuzzy.LogLevel) fuzzy.Logger {
	l := logrus.New()
	switch lvl {
	case fuzzy.LogPanic:
		l.Level = logrus.PanicLevel
	case fuzzy.LogFatal:
		l.Level = logrus.FatalLevel
	case fuzzy.LogError:
		l.Level = logrus.ErrorLevel
	case fuzzy.LogWarn:
		l.Level = logrus.WarnLevel
	case fuzzy.LogInfo:
		l.Level = logrus.InfoLevel
	case fuzzy.LogDebug:
		l.Level = logrus.DebugLevel
	}

	l.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
	return logrusLogger{l}
}
