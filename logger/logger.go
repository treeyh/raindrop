package logger

import (
	"context"
	"github.com/treeyh/raindrop/model"
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
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Fatal(context.Context, string, ...interface{})
}

func New(writer IWriter, config model.RainDropLogConfig) ILogger {
	var (
		debugStr = "%s\n[debug] "
		infoStr  = "%s\n[info] "
		warnStr  = "%s\n[warn] "
		errStr   = "%s\n[error] "
		fatalStr = "%s\n[error] "
	)

	if config.Colorful {
		debugStr = Green + "%s\n" + Reset + Green + "[debug] " + Reset
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		fatalStr = Red + "%s\n" + Reset + RedBold + "[fatal] " + Reset
	}

	return &_logger{
		IWriter:           writer,
		RainDropLogConfig: config,
		debugStr:          debugStr,
		infoStr:           infoStr,
		warnStr:           warnStr,
		errStr:            errStr,
		fatalStr:          fatalStr,
	}
}

type _logger struct {
	IWriter
	model.RainDropLogConfig
	debugStr, infoStr, warnStr, errStr, fatalStr string
}

// LogMode log mode
func (l *_logger) LogMode(level LogLevel) ILogger {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Debug print info
func (l _logger) Debug(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Debug {
		l.Printf(ctx, l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Info print info
func (l _logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Info {
		l.Printf(ctx, l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l _logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Warn {
		l.Printf(ctx, l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l _logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Error {
		l.Printf(ctx, l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Fatal print error messages
func (l _logger) Fatal(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Fatal {
		l.Printf(ctx, l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
	panic(msg)
}
