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

// Levels 日志的级别,字符串表示
var Levels = map[string]int{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
	"FATAL": FATAL,
}

// GetLevel 返回日志级别,没有就返回DEBUG
func GetLevel(level string) int {
	l := Levels[strings.ToUpper(level)]
	if l == 0 {
		return DEBUG
	}
	return l
}

// Logger 表示有日志级别的
type Logger struct {
	level int
	w     io.Writer
	*log.Logger
}

var stdLog = New(INFO, os.Stderr, log.LstdFlags|log.Lshortfile)

// Init 重新初始化
func Init(level int, l io.Writer, flag int) {
	stdLog = &Logger{
		level:  level,
		w:      l,
		Logger: log.New(l, "", flag),
	}
}

// New 返回*Logger
func New(level int, l io.Writer, flag int) *Logger {
	return &Logger{
		level:  level,
		w:      l,
		Logger: log.New(l, "", flag),
	}
}

// NewCached 返回日志缓存在内存中
func NewCached(level int, flag int) *Logger {
	l := &bytes.Buffer{}
	return &Logger{
		level:  level,
		w:      l,
		Logger: log.New(l, "", flag),
	}
}

// Log 表示记录相应等级以上的日志
func (l *Logger) Log(level int, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// Debugf 表示记录DEBUG以上日志
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// Infof 表示记录INFO以上日志
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// Warnf 表示记录WARN以上日志
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level > WARN {
		return
	}
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// Errorf 表示记录ERROR以上日志
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// GetOutput 返回日志输出到的地方
func (l *Logger) GetOutput() io.Writer {
	return l.w
}

// Write 记录日志
func (l *Logger) Write(format string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// Write 记录日志
func Write(format string, v ...interface{}) {
	_ = stdLog.Output(2, fmt.Sprintf(format, v...))
}

// Log 表示记录相应等级以上的日志
func Log(level int, format string, v ...interface{}) {
	if level < stdLog.level {
		return
	}
	_ = stdLog.Output(2, fmt.Sprintf(format, v...))
}

// Fatalf 表示记录日志然后退出
func Fatalf(format string, v ...interface{}) {
	_ = stdLog.Output(2, fmt.Sprintf("Fatal: "+format, v...))
	os.Exit(1)
}

// Debugf 表示记录DEBUG以上日志
func Debugf(format string, v ...interface{}) {
	if stdLog.level > DEBUG {
		return
	}
	_ = stdLog.Output(2, fmt.Sprintf("Debug: "+format, v...))
}

// Infof 表示记录INFO以上日志
func Infof(format string, v ...interface{}) {
	if stdLog.level > INFO {
		return
	}
	_ = stdLog.Output(2, fmt.Sprintf("Info: "+format, v...))
}

// Warnf 表示记录WARN以上日志
func Warnf(format string, v ...interface{}) {
	if stdLog.level > WARN {
		return
	}
	_ = stdLog.Output(2, fmt.Sprintf("Warn: "+format, v...))
}

// Errorf 表示记录ERROR以上日志
func Errorf(format string, v ...interface{}) {
	if stdLog.level > ERROR {
		return
	}
	_ = stdLog.Output(2, fmt.Sprintf("Error: "+format, v...))
}
