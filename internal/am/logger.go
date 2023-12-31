package am

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	ErrorLevel
)

type Logger interface {
	SetLogLevel(level LogLevel)
	Debug(v ...any)
	Debugf(format string, a ...any)
	Info(v ...any)
	Infof(format string, a ...any)
	Error(v ...any)
	Errorf(format string, a ...any)
}

type SimpleLogger struct {
	debug    *log.Logger
	info     *log.Logger
	error    *log.Logger
	logLevel LogLevel
}

func NewLogger(logLevel string) *SimpleLogger {
	level := ToValidLevel(logLevel)
	return &SimpleLogger{
		debug:    log.New(os.Stdout, "[DBG] ", log.LstdFlags),
		info:     log.New(os.Stdout, "[INF] ", log.LstdFlags),
		error:    log.New(os.Stderr, "[ERR] ", log.LstdFlags),
		logLevel: level,
	}
}

func (l *SimpleLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

func (l *SimpleLogger) Debug(v ...any) {
	if l.logLevel <= DebugLevel {
		l.debug.Println(v...)
	}
}

func (l *SimpleLogger) Debugf(format string, a ...any) {
	if l.logLevel <= DebugLevel {
		message := fmt.Sprintf(format, a...)
		l.debug.Println(message)
	}
}

func (l *SimpleLogger) Info(v ...any) {
	if l.logLevel <= InfoLevel {
		l.info.Println(v...)
	}
}

func (l *SimpleLogger) Infof(format string, a ...any) {
	if l.logLevel <= InfoLevel {
		message := fmt.Sprintf(format, a...)
		l.info.Println(message)
	}
}

func (l *SimpleLogger) Error(v ...interface{}) {
	if l.logLevel <= ErrorLevel {
		message := fmt.Sprint(v...)
		l.error.Println(message)
	}
}

func (l *SimpleLogger) Errorf(format string, a ...interface{}) {
	if l.logLevel <= ErrorLevel {
		message := fmt.Sprintf(format, a...)
		l.error.Println(message)
	}
}

func ToValidLevel(level string) LogLevel {
	level = strings.ToLower(level)

	switch level {
	case "debug", "dbg":
		return DebugLevel
	case "info", "inf":
		return InfoLevel
	case "error", "err":
		return ErrorLevel
	default:
		return ErrorLevel
	}
}

// SetDebugOutput set the internal logger.
// Used for package testing.
func (sl *SimpleLogger) SetDebugOutput(debug *bytes.Buffer) {
	sl.debug = log.New(debug, "", 1)
}

// SetInfoOutput set the internal logger.
// Used for package testing.
func (sl *SimpleLogger) SetInfoOutput(info *bytes.Buffer) {
	sl.info = log.New(info, "", 1)
}

// SetErrorOutput set the internal logger.
// Used for package testing.
func (sl *SimpleLogger) SetErrorOutput(error *bytes.Buffer) {
	sl.error = log.New(error, "", 1)
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[1] = unicode.ToUpper(runes[0])
	return string(runes)
}
