package mylog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
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

func getLevelString(level int) string {
	for k, v := range Levels {
		if v == level {
			return k
		}
	}
	return ""
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
	mu          sync.Mutex
	level       int
	levelString string
	flag        int
	json        bool
	w           io.Writer
	*log.Logger

	preM map[string]interface{}
}

var stdLog = New(INFO, os.Stderr, log.LstdFlags|log.Lshortfile, false)

// Init 重新初始化
func Init(level int, l io.Writer, flag int, json bool) {
	stdLog = New(level, l, flag, json)
}

// New 返回*Logger
func New(level int, l io.Writer, flag int, json bool) *Logger {
	return &Logger{
		level:       level,
		levelString: getLevelString(level),
		flag:        flag,
		json:        json,
		w:           l,
		Logger:      log.New(l, "", flag),
		preM:        make(map[string]interface{}),
	}
}

// NewCached 返回日志缓存在内存中
func NewCached(level int, flag int, json bool) *Logger {
	l := &bytes.Buffer{}
	return New(level, l, flag, json)
}

// Predefine 表示记录额外的字段
func (l *Logger) Predefine(m map[string]interface{}) {
	l.mu.Lock()
	l.preM = m
	l.mu.Unlock()
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

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) formatHeader(t time.Time, file string, line int) map[string]interface{} {
	m := make(map[string]interface{})
	if l.flag&log.LUTC != 0 {
		t = t.UTC()
	}

	var buf []byte
	if l.flag&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.flag&log.Ldate != 0 {
			year, month, day := t.Date()
			itoa(&buf, year, 4)
			buf = append(buf, '/')
			itoa(&buf, int(month), 2)
			buf = append(buf, '/')
			itoa(&buf, day, 2)
			buf = append(buf, ' ')
		}
		if l.flag&(log.Ltime|log.Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(&buf, hour, 2)
			buf = append(buf, ':')
			itoa(&buf, min, 2)
			buf = append(buf, ':')
			itoa(&buf, sec, 2)
			if l.flag&log.Lmicroseconds != 0 {
				buf = append(buf, '.')
				itoa(&buf, t.Nanosecond()/1e3, 6)
			}
		}
	}
	m["logtime"] = string(buf)

	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		if l.flag&log.Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		m["filename"] = file
		m["lineno"] = line
	}

	return m
}

// Output 格式化日志
func (l *Logger) Output(calldepth int, s string) error {
	if !l.json {
		return l.Logger.Output(calldepth, s)
	}

	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	m := l.formatHeader(now, file, line)
	m["log_level"] = l.levelString
	m["msg"] = s

	for k, v := range l.preM {
		m[k] = v
	}

	b, err := json.Marshal(m)

	if err != nil {
		panic("json.Marshal error")
	}

	b = append(b, '\n')

	_, err = l.w.Write(b)
	return err
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

// Predefine 表示记录额外的字段
func Predefine(m map[string]interface{}) {
	stdLog.Predefine(m)
}
