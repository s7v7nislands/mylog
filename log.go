package mylog

import (
	"fmt"
	"log"
	"strings"
)

const (
	DEBUG = iota << 1
	INFO
	WARN
	ERROR
)

var Levels = map[string]int{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
}

func GetLevel(level string) int {
	l := Levels[strings.ToUpper(level)]
	if l == 0 {
		return DEBUG
	}
	return l
}

type Logger struct {
	level int
	*log.Logger
}

func New(level int, log *log.Logger) *Logger {
	return &Logger{level, log}
}

func (l *Logger) Log(level int, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level > WARN {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Err(format string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}
