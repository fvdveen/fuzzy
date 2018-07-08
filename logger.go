package fuzzy

import (
	"fmt"
	"log"
	"os"
)

const (
	// LogXxx are the different log levels
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
	LogFatal
	LogPanic
)

// LogLevel is the level for the logger
type LogLevel int

// LogLevelFrom gets the log level from the string
// it returns 0 (LogDebug) when the string is unrecognized
func LogLevelFrom(lvl string) LogLevel {
	switch lvl {
	case "debug":
		return LogDebug
	case "info":
		return LogInfo
	case "warn", "warning":
		return LogWarn
	case "error":
		return LogError
	case "fatal":
		return LogFatal
	case "panic":
		return LogPanic
	default:
		return 0
	}
}

// LoggerGenerator creates a Logger
type LoggerGenerator func(LogLevel) Logger

// Logger is the logger used throughout the bot
type Logger interface {
	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
}

type defaultLogger struct {
	fields map[string]interface{}
	level  LogLevel
	logger *log.Logger
}

// DefaultLogger is the default logger
func DefaultLogger(lvl LogLevel) Logger {
	l := &defaultLogger{
		fields: make(map[string]interface{}),
		level:  lvl,
	}

	l.logger = log.New(os.Stderr, "", log.LstdFlags)

	return l
}

func (l *defaultLogger) WithField(k string, v interface{}) Logger {
	l2 := DefaultLogger(l.level).(*defaultLogger)
	for k2, v2 := range l.fields {
		l2.fields[k2] = v2
	}

	l2.fields[k] = v

	return l2
}

func (l *defaultLogger) WithFields(fs map[string]interface{}) Logger {
	l2 := DefaultLogger(l.level).(*defaultLogger)
	for k, v := range l.fields {
		l2.fields[k] = v
	}

	for k, v := range fs {
		l2.fields[k] = v
	}

	return l2
}

func (l *defaultLogger) Debugf(format string, a ...interface{}) {
	if l.level > LogDebug {
		return
	}
	format = "[]: " + format
	format = fmt.Sprintf(format, a...)
	for k, v := range l.fields {
		format = fmt.Sprintf("%s [%s]=[%v]", format, k, v)
	}
	l.logger.Println(format)
}

func (l *defaultLogger) Infof(format string, a ...interface{}) {
	if l.level > LogInfo {
		return
	}
	format = "[]: " + format
	format = fmt.Sprintf(format, a...)
	for k, v := range l.fields {
		format = fmt.Sprintf("%s [%s]=[%v]", format, k, v)
	}
	l.logger.Println(format)
}

func (l *defaultLogger) Printf(format string, a ...interface{}) {
	format = "[]: " + format
	format = fmt.Sprintf(format, a...)
	for k, v := range l.fields {
		format = fmt.Sprintf("%s [%s]=[%v]", format, k, v)
	}
	l.logger.Println(format)
}

func (l *defaultLogger) Warnf(format string, a ...interface{}) {
	if l.level > LogWarn {
		return
	}
	format = "[]: " + format
	format = fmt.Sprintf(format, a...)
	for k, v := range l.fields {
		format = fmt.Sprintf("%s [%s]=[%v]", format, k, v)
	}
	l.logger.Println(format)
}

func (l *defaultLogger) Errorf(format string, a ...interface{}) {
	if l.level > LogError {
		return
	}
	format = "[]: " + format
	format = fmt.Sprintf(format, a...)
	for k, v := range l.fields {
		format = fmt.Sprintf("%s [%s]=[%v]", format, k, v)
	}
	l.logger.Println(format)
}

func (l *defaultLogger) Fatalf(format string, a ...interface{}) {
	if l.level > LogFatal {
		os.Exit(1)
		return
	}
	format = "[]: " + format
	format = fmt.Sprintf(format, a...)
	for k, v := range l.fields {
		format = fmt.Sprintf("%s [%s]=[%v]", format, k, v)
	}
	l.logger.Fatalln(format)
}

func (l *defaultLogger) Debug(a ...interface{}) {
	if l.level > LogDebug {
		return
	}
	l.logger.Println("[DEBUG]: " + fmt.Sprintln(a...))
}

func (l *defaultLogger) Info(a ...interface{}) {
	if l.level > LogInfo {
		return
	}
	l.logger.Println("[INFO]: " + fmt.Sprintln(a...))
}

func (l *defaultLogger) Warn(a ...interface{}) {
	if l.level > LogWarn {
		return
	}
	l.logger.Println("[WARN]: " + fmt.Sprintln(a...))
}

func (l *defaultLogger) Error(a ...interface{}) {
	if l.level > LogError {
		return
	}
	l.logger.Println("[ERROR]: " + fmt.Sprintln(a...))
}

func (l *defaultLogger) Fatal(a ...interface{}) {
	if l.level > LogFatal {
		os.Exit(1)
		return
	}
	l.logger.Fatalln("[FATAL]: " + fmt.Sprintln(a...))
}

func (l *defaultLogger) Panic(a ...interface{}) {
	panic("[PANIC]: " + fmt.Sprintln(a...))
}
