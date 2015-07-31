package mylog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// debug的级别
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

// debug的级别,字符串表示
var Levels = map[string]int{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
	"FATAL": FATAL,
}

// GetLevel返回日志级别,没有就返回DEBUG
func GetLevel(level string) int {
	l := Levels[strings.ToUpper(level)]
	if l == 0 {
		return DEBUG
	}
	return l
}

// Logger表示有日志级别的
type logger struct {
	level int
	l     io.Writer
	*log.Logger
}

var stdLog = New(INFO, os.Stderr, log.LstdFlags|log.Lshortfile)

// New返回*Logger
func New(level int, l io.Writer, flag int) *logger {
	return &logger{level: level, l: l, Logger: log.New(l, "", flag)}
}

// NewCached返回日志缓存在内存中
func NewCached(level int, flag int) *logger {
	l := &bytes.Buffer{}
	return &logger{level: level, l: l, Logger: log.New(l, "", flag)}
}

func (l *logger) Log(level int, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) Debug(format string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) Info(format string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) Warn(format string, v ...interface{}) {
	if l.level > WARN {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) Error(format string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) GetOutput() io.Writer {
	return l.l
}

func (l *logger) Write(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

func Write(format string, v ...interface{}) {
	stdLog.Output(2, fmt.Sprintf(format, v...))
}

func Log(level int, format string, v ...interface{}) {
	if level < stdLog.level {
		return
	}
	stdLog.Output(2, fmt.Sprintf(format, v...))
}

func Fatal(format string, v ...interface{}) {
	stdLog.Output(2, fmt.Sprintf("Fatal: "+format, v...))
	os.Exit(1)
}

func Debug(format string, v ...interface{}) {
	if stdLog.level > DEBUG {
		return
	}
	stdLog.Output(2, fmt.Sprintf("Debug: "+format, v...))
}

func Info(format string, v ...interface{}) {
	if stdLog.level > INFO {
		return
	}
	stdLog.Output(2, fmt.Sprintf("Info: "+format, v...))
}

func Warn(format string, v ...interface{}) {
	if stdLog.level > WARN {
		return
	}
	stdLog.Output(2, fmt.Sprintf("Warn: "+format, v...))
}

func Error(format string, v ...interface{}) {
	if stdLog.level > ERROR {
		return
	}
	stdLog.Output(2, fmt.Sprintf("Error: "+format, v...))
}
