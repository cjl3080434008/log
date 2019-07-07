package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var DefaultGlogLogger = GlogLogger{InfoLevel, false, &Writer{}}

type GlogLogger struct {
	Level     Level
	ANSIColor bool
	Writer    io.Writer
}

type GlogEvent struct {
	buf       []byte
	level     Level
	Level     Level
	ANSIColor bool
	Write     func(p []byte) (n int, err error)
}

func (l GlogLogger) Info(args ...interface{}) {
	l.WithLevel(InfoLevel).Print(args...)
}

func (l GlogLogger) Infoln(args ...interface{}) {
	l.WithLevel(InfoLevel).Println(args...)
}

func (l GlogLogger) Infof(format string, args ...interface{}) {
	l.WithLevel(InfoLevel).Printf(format, args...)
}

func (l GlogLogger) InfoDepth(depth int, args ...interface{}) {
	l.WithLevel(InfoLevel).PrintDepth(depth, args...)
}

func (l GlogLogger) Warning(args ...interface{}) {
	l.WithLevel(WarnLevel).Print(args...)
}

func (l GlogLogger) Warningln(args ...interface{}) {
	l.WithLevel(WarnLevel).Println(args...)
}

func (l GlogLogger) Warningf(format string, args ...interface{}) {
	l.WithLevel(WarnLevel).Printf(format, args...)
}

func (l GlogLogger) WarningDepth(depth int, args ...interface{}) {
	l.WithLevel(WarnLevel).PrintDepth(depth, args...)
}

func (l GlogLogger) Error(args ...interface{}) {
	l.WithLevel(ErrorLevel).Print(args...)
}

func (l GlogLogger) Errorln(args ...interface{}) {
	l.WithLevel(ErrorLevel).Println(args...)
}

func (l GlogLogger) Errorf(format string, args ...interface{}) {
	l.WithLevel(ErrorLevel).Printf(format, args...)
}

func (l GlogLogger) ErrorDepth(depth int, args ...interface{}) {
	l.WithLevel(ErrorLevel).PrintDepth(depth, args...)
}

func (l GlogLogger) Fatal(args ...interface{}) {
	l.WithLevel(FatalLevel).Print(args...)
}

func (l GlogLogger) Fatalln(args ...interface{}) {
	l.WithLevel(FatalLevel).Println(args...)
}

func (l GlogLogger) Fatalf(format string, args ...interface{}) {
	l.WithLevel(FatalLevel).Printf(format, args...)
}

func (l GlogLogger) FatalDepth(depth int, args ...interface{}) {
	l.WithLevel(FatalLevel).PrintDepth(depth, args...)
}

func (l GlogLogger) V(level int) bool {
	return level >= int(l.Level)
}

var gepool = sync.Pool{
	New: func() interface{} {
		return new(GlogEvent)
	},
}

var pid = int64(os.Getpid())

func (l GlogLogger) WithLevel(level Level) (e *GlogEvent) {
	if level < l.Level {
		return
	}
	// [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg
	e = gepool.Get().(*GlogEvent)
	e.buf = e.buf[:0]
	e.level = level
	e.Level = l.Level
	e.ANSIColor = l.ANSIColor
	e.Write = l.Writer.Write
	// level
	switch level {
	case DebugLevel:
		e.colorize('D', ColorGreen)
	case InfoLevel:
		e.colorize('I', ColorCyan)
	case WarnLevel:
		e.colorize('W', ColorYellow)
	case ErrorLevel:
		e.colorize('E', ColorRed)
	case FatalLevel:
		e.colorize('F', ColorRed)
	default:
		e.colorize('?', ColorRed)
	}
	// time
	now := timeNow()
	if e.ANSIColor {
		e.buf = append(e.buf, ColorDarkGray...)
	}
	e.time(now)
	if e.ANSIColor {
		e.buf = append(e.buf, ColorReset...)
	}
	e.buf = append(e.buf, ' ')
	// threadid
	e.buf = strconv.AppendInt(e.buf, pid, 10)
	e.buf = append(e.buf, ' ')

	return
}

func (e *GlogEvent) Printf(format string, args ...interface{}) {
	if e == nil {
		return
	}
	e.caller(3)
	msg := fmt.Sprintf(format, args...)
	e.buf = append(e.buf, msg...)
	e.buf = append(e.buf, '\n')
	e.Write(e.buf)
	gepool.Put(e)
	if e.level >= FatalLevel {
		panic(msg)
	}
}

func (e *GlogEvent) Print(args ...interface{}) {
	if e == nil {
		return
	}
	e.caller(3)
	msg := fmt.Sprint(args...)
	e.buf = append(e.buf, msg...)
	e.buf = append(e.buf, '\n')
	e.Write(e.buf)
	gepool.Put(e)
	if e.level >= FatalLevel {
		panic(msg)
	}
}

func (e *GlogEvent) Println(args ...interface{}) {
	if e == nil {
		return
	}
	e.caller(3)
	msg := fmt.Sprintln(args...)
	e.buf = append(e.buf, msg...)
	e.buf = append(e.buf, '\n')
	e.Write(e.buf)
	gepool.Put(e)
	if e.level >= FatalLevel {
		panic(msg)
	}
}

func (e *GlogEvent) PrintDepth(depth int, args ...interface{}) {
	if e == nil {
		return
	}
	e.caller(3 + depth)
	msg := fmt.Sprint(args...)
	e.buf = append(e.buf, msg...)
	e.buf = append(e.buf, '\n')
	e.Write(e.buf)
	gepool.Put(e)
	if e.level >= FatalLevel {
		panic(msg)
	}
}

func (e *GlogEvent) time(now time.Time) {
	var n = len(e.buf)
	e.buf = append(e.buf, "0102 15:04:05.999999"...)
	var a, b int
	// month
	a = int(now.Month())
	b = a / 10
	e.buf[n+1] = byte('0' + a - 10*b)
	e.buf[n] = byte('0' + b)
	// day
	a = now.Day()
	b = a / 10
	e.buf[n+3] = byte('0' + a - 10*b)
	e.buf[n+2] = byte('0' + b)
	// hour
	a = now.Hour()
	b = a / 10
	e.buf[n+6] = byte('0' + a - 10*b)
	e.buf[n+5] = byte('0' + b)
	// minute
	a = now.Minute()
	b = a / 10
	e.buf[n+9] = byte('0' + a - 10*b)
	e.buf[n+8] = byte('0' + b)
	// second
	a = now.Second()
	b = a / 10
	e.buf[n+12] = byte('0' + a - 10*b)
	e.buf[n+11] = byte('0' + b)
	// milli second
	a = now.Nanosecond() / 1000
	b = a / 10
	e.buf[n+19] = byte('0' + a - 10*b)
	a = b
	b = a / 10
	e.buf[n+18] = byte('0' + a - 10*b)
	a = b
	b = a / 10
	e.buf[n+17] = byte('0' + a - 10*b)
	a = b
	b = a / 10
	e.buf[n+16] = byte('0' + a - 10*b)
	a = b
	b = a / 10
	e.buf[n+15] = byte('0' + a - 10*b)
	e.buf[n+14] = byte('0' + b)
}

func (e *GlogEvent) caller(depth int) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	if line < 0 {
		line = 0
	}
	e.buf = append(e.buf, file...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendInt(e.buf, int64(line), 10)
	e.colorize(']', ColorCyan)
	e.buf = append(e.buf, ' ')
}

func (e *GlogEvent) colorize(b byte, c ANSIColor) {
	if e == nil {
		return
	}

	if !e.ANSIColor {
		e.buf = append(e.buf, b)
		return
	}

	e.buf = append(e.buf, c...)
	e.buf = append(e.buf, b)
	e.buf = append(e.buf, ColorReset...)
}