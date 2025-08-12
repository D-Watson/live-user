package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	JsonEncoder       = "json"
	ConsoleEncoder    = "console"
	DefaultCallerSkip = 1
	DefaultLogFile    = "./log/app.log"
	DefaultErrLogFile = "./log/err.log"
	DefaultMaxSize    = 100 * 1024 * 1024
	DefaultMaxAge     = 7
	DefaultBackups    = 100
	DefaultTimeLayout = "2006-01-02T15:04:05.000"

	TraceIDFlag = "traceid"
	SpanIDFlag  = "spanid"
)

var globallogger *zap.Logger
var openTrace bool

type Config struct {
	File     string `yaml:"file" json:"file"`
	ErrFile  string `yaml:"errfile" json:"errfile"`
	Encode   string `yaml:"encode" json:"encode"`   // json or console
	Level    string `yaml:"level" json:"level"`     // debug info warn error panic fatal
	MaxSize  int64  `yaml:"maxsize" json:"maxsize"` // rotate is file need set max size	MB
	MaxDay   int    `yaml:"maxday" json:"maxday"`   // rotate days
	Backups  int    `yaml:"backups" json:"backups"`
	Compress bool   `yaml:"compress" json:"compress"`
	Caller   int    `yaml:"caller" json:"caller"`
	Trace    bool   `yaml:"trace" json:"trace"` // wither log trace info
}

func (cfg Config) parseTimeHook(file string) *rotatelogs.RotateLogs {
	newHook, err := rotatelogs.New(
		file+".%Y%m%d%H",
		rotatelogs.WithLinkName(file),
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxDay)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithRotationSize(cfg.MaxSize*1024*1024),
	)

	if err != nil {
		fmt.Println("time rotate failed", err)
		return defaultTimeRotateHook(cfg.File)
	}
	return newHook
}

func init() {
	openTrace = true
	logHook := defaultFileRotateHook(DefaultLogFile)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(DefaultTimeLayout)
	encode := NewCustomEncoder(encoderConfig)
	core := zapcore.NewCore(
		encode,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logHook)),
		zap.NewAtomicLevel())

	caller := zap.AddCaller()
	callerskip := zap.AddCallerSkip(DefaultCallerSkip)
	develop := zap.Development()

	globallogger = zap.New(core, caller, callerskip, develop)
}

func defaultFileRotateHook(file string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    DefaultMaxSize,
		MaxAge:     DefaultMaxAge,
		MaxBackups: DefaultBackups,
		LocalTime:  true,
		Compress:   false,
	}
}

func defaultTimeRotateHook(file string) *rotatelogs.RotateLogs {
	hook, _ := rotatelogs.New(
		file+"%Y%m%d%H",
		rotatelogs.WithLinkName(file),
		rotatelogs.WithMaxAge(3*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithRotationSize(DefaultMaxSize))

	return hook
}

func Init(cfg *Config) {
	var logHook io.Writer
	var errLogHook io.Writer

	if cfg.File == "" {
		cfg.File = DefaultLogFile
	}
	logHook = cfg.parseTimeHook(cfg.File)

	if cfg.ErrFile == "" {
		cfg.ErrFile = DefaultErrLogFile
	}
	errLogHook = cfg.parseTimeHook(cfg.ErrFile)

	openTrace = cfg.Trace

	encoderconf := zap.NewProductionEncoderConfig()
	encoderconf.EncodeTime = zapcore.TimeEncoderOfLayout(DefaultTimeLayout)

	encode := zapcore.NewJSONEncoder(encoderconf)
	if cfg.Encode == ConsoleEncoder {
		encode = NewCustomEncoder(encoderconf)
	}

	loglevel, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		fmt.Println("level illegal, use default level: info")
		panic(err)
	}

	core := zapcore.NewTee(
		zapcore.NewCore( // logCore
			encode,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logHook)),
			loglevel),
		zapcore.NewCore( // errLogCore
			encode,
			zapcore.AddSync(errLogHook),
			zap.NewAtomicLevelAt(zap.ErrorLevel)),
	)

	caller := zap.AddCaller()
	callerskip := zap.AddCallerSkip(DefaultCallerSkip)
	if cfg.Caller != 0 {
		callerskip = zap.AddCallerSkip(cfg.Caller)
	}

	develop := zap.Development()

	globallogger = zap.New(core, caller, callerskip, develop)

	InitGinLogs(cfg.File)
}

func InitGinLogs(file string) {
	//f := &logs.PatternLogFormatter{
	//	Pattern:    "%w|%T|%f:%n|beego framework log|msg=%m\n ",
	//	WhenFormat: DefaultTimeLayout,
	//}
	//logs.RegisterFormatter("xiaomi-pattern", f)

	//conf := `{"filename": "` + file + `.beegoframework","maxsize": 10000000,"formatter": "xiaomi-pattern","perm":"0777"}`
	//_ = logs.SetLogger(logs.AdapterConsole, conf) 	控制台输出

}

func Debug(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	globallogger.Debug(msg, getLogTrace(ctx)...)
}

func Info(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	globallogger.Info(msg, getLogTrace(ctx)...)
}

func Warn(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	globallogger.Warn(msg, getLogTrace(ctx)...)
}

func Error(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	globallogger.Error(msg, getLogTrace(ctx)...)
}

func Panic(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	globallogger.Panic(msg, getLogTrace(ctx)...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	globallogger.Fatal(msg, getLogTrace(ctx)...)
}

func Debugf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	globallogger.Debug(msg, getLogTrace(ctx)...)
}

func Infof(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	globallogger.Info(msg, getLogTrace(ctx)...)
}

func Warnf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	globallogger.Warn(msg, getLogTrace(ctx)...)
}

func Errorf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	globallogger.Error(msg, getLogTrace(ctx)...)
}

func Panicf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	globallogger.Panic(msg, getLogTrace(ctx)...)
}

func Fatalf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	globallogger.Fatal(msg, getLogTrace(ctx)...)
}

func getLogTrace(ctx context.Context) []zap.Field {
	if !openTrace {
		return nil
	}

	return []zap.Field{
		//zap.Any(TraceIDFlag, trace.TraceID(ctx)),
		//zap.Any(SpanIDFlag, trace.SpanID(ctx)),
	}
}
