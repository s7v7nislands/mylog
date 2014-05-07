package mylog

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// debug的级别
const (
	DEBUG = iota << 1
	INFO
	WARN
	ERROR
)

// debug的级别,字符串表示
var Levels = map[string]int{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
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
	*log.Logger
}

// loggerGroup表示可以输出到多个日志
type LoggerGroup struct {
	g []*logger
}

var stdLog = newLogger(INFO, log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))
var stdLogGroup = New(stdLog)

// New返回*Logger
func newLogger(level int, log *log.Logger) *logger {
	return &logger{level, log}
}

// NewGroup返回*LoggerGroup
func New(l ...*logger) *LoggerGroup {
	g := &LoggerGroup{}
	g.g = append(g.g, l...)

	return g
}

func Init(level int, log *log.Logger) {
	stdLogGroup.Init(level, log)
}

func AddHandler(level int, log *log.Logger) {
	l := newLogger(level, log)
	stdLogGroup.g = append(stdLogGroup.g, l)
}

func (g *LoggerGroup) Init(level int, log *log.Logger) {
	l := newLogger(level, log)
	g.g = []*logger{l}
}

func (g *LoggerGroup) AddHandler(level int, log *log.Logger) {
	l := newLogger(level, log)
	g.g = append(g.g, l)
}

func (l *logger) log(level int, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) debug(format string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) info(format string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) warn(format string, v ...interface{}) {
	if l.level > WARN {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) err(format string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func (g *LoggerGroup) Log(level int, format string, v ...interface{}) {
	for _, l := range g.g {
		if level < l.level {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (g *LoggerGroup) Debug(format string, v ...interface{}) {
	for _, l := range g.g {
		if l.level > DEBUG {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (g *LoggerGroup) Info(format string, v ...interface{}) {
	for _, l := range g.g {
		if l.level > INFO {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (g *LoggerGroup) Warn(format string, v ...interface{}) {
	for _, l := range g.g {
		if l.level > WARN {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (g *LoggerGroup) Err(format string, v ...interface{}) {
	for _, l := range g.g {
		if l.level > ERROR {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func Log(level int, format string, v ...interface{}) {
	for _, l := range stdLogGroup.g {
		if level < l.level {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func Fatal(format string, v ...interface{}) {
	for _, l := range stdLogGroup.g {
		if l.level > ERROR {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func Debug(format string, v ...interface{}) {
	for _, l := range stdLogGroup.g {
		if l.level > DEBUG {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	for _, l := range stdLogGroup.g {
		if l.level > INFO {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func Warn(format string, v ...interface{}) {
	for _, l := range stdLogGroup.g {
		if l.level > WARN {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func Err(format string, v ...interface{}) {
	for _, l := range stdLogGroup.g {
		if l.level > ERROR {
			continue
		}
		l.Output(2, fmt.Sprintf(format, v...))
	}
}
