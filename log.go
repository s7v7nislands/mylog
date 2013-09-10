package mylog

import (
	"log"
)

const (
	DEBUG = iota << 1
	INFO
	WARN
	ERROR
)

type Logger struct {
	level int
	*log.Logger
}

func New(level int, log *log.Logger) *Logger {
	return &Logger{level, log}
}

func (l *Logger) Log(level int, fmt string, v ...interface{}) {
	if level < l.level {
		return
	}
	l.Printf(fmt, v...)
}

func (l *Logger) Debug(fmt string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	l.Printf(fmt, v...)
}

func (l *Logger) Info(fmt string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	l.Printf(fmt, v...)
}

func (l *Logger) Warn(fmt string, v ...interface{}) {
	if l.level > WARN {
		return
	}
	l.Printf(fmt, v...)
}

func (l *Logger) Error(fmt string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	l.Printf(fmt, v...)
}
