package logger

import (
	"context"
	"fmt"
	"github.com/treeyh/raindrop/utils"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LogLevel 日志级别
type LogLevel int

const (
	Debug LogLevel = iota + 1
	Info
	Warn
	Error
	Fatal
)

// IWriter 写日志方法
type IWriter interface {
	Printf(context.Context, string, ...interface{})
}

type ILogger interface {
	LogMode(LogLevel) ILogger
	GetLogLevel() LogLevel
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Fatal(context.Context, string, ...interface{})
}

type DefaultWriter struct {
}

func (dw *DefaultWriter) Printf(ctx context.Context, msg string, data ...interface{}) {
	fmt.Println(append([]interface{}{msg}, data...))
}

func New(writer IWriter, logLevel LogLevel, colorful bool) ILogger {
	var (
		debugStr = "[debug] %s"
		infoStr  = "[info] %s"
		warnStr  = "[warn] %s"
		errStr   = "[error] %s"
		fatalStr = "[fatal] %s"
	)

	if colorful {
		debugStr = Green + "[debug] " + "%s" + Reset + Green + Reset
		infoStr = Green + "[info] " + "%s" + Reset + Green + Reset
		warnStr = BlueBold + "[warn] " + "%s" + Reset + Magenta + Reset
		errStr = Magenta + "[error] " + "%s" + Reset + Red + Reset
		fatalStr = Red + "[fatal] " + "%s" + Reset + RedBold + Reset
	}

	return &_logger{
		IWriter:  writer,
		LogLevel: logLevel,
		Colorful: colorful,
		debugStr: debugStr,
		infoStr:  infoStr,
		warnStr:  warnStr,
		errStr:   errStr,
		fatalStr: fatalStr,
	}
}

func NewDefault() ILogger {
	d := DefaultWriter{}
	return New(&d, Info, true)
}

type _logger struct {
	IWriter
	LogLevel                                     LogLevel
	Colorful                                     bool
	debugStr, infoStr, warnStr, errStr, fatalStr string
}

// LogMode log mode
func (l *_logger) LogMode(level LogLevel) ILogger {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// GetLogLevel 获取日志级别
func (l *_logger) GetLogLevel() LogLevel {
	return l.LogLevel
}

// Debug print info
func (l _logger) Debug(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel <= Debug {
		l.Printf(ctx, fmt.Sprintf(l.debugStr, msg), append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Info print info
func (l _logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel <= Info {
		l.Printf(ctx, fmt.Sprintf(l.infoStr, msg), append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l _logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel <= Warn {
		l.Printf(ctx, fmt.Sprintf(l.warnStr, msg), append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l _logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel <= Error {
		l.Printf(ctx, fmt.Sprintf(l.errStr, msg), append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Fatal print error messages
func (l _logger) Fatal(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel <= Fatal {
		l.Printf(ctx, fmt.Sprintf(l.fatalStr, msg), append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
	panic(msg)
}
