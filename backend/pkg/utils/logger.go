package utils

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerInterface interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Logger struct {
	*zap.SugaredLogger
}

var LoggerInstance *Logger
var once sync.Once

func NewLogger(level string) (*Logger, error) {
	var err error

	once.Do(func() {
		var zapLevel zapcore.Level

		zapLevel, e := zapcore.ParseLevel(level)
		if e != nil {
			err = fmt.Errorf("invalid logger level: %w", e)
			return
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:          "time",
			LevelKey:         "level",
			NameKey:          "logger",
			CallerKey:        "caller",
			MessageKey:       "msg",
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeLevel:      customColorLevelEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
		}

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zapLevel),
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

		LoggerInstance = &Logger{
			zapLogger.Sugar(),
		}
	})

	return LoggerInstance, err
}

func customColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString("[\033[36mDEBUG\033[0m]")
	case zapcore.InfoLevel:
		enc.AppendString("[\033[32mINFO\033[0m]")
	case zapcore.WarnLevel:
		enc.AppendString("[\033[33mWARN\033[0m]")
	case zapcore.ErrorLevel:
		enc.AppendString("[\033[31mERROR\033[0m]")
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString("[\033[35mFATAL\033[0m]")
	default:
		enc.AppendString("[UNKNOWN]")
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.SugaredLogger.Debug(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.SugaredLogger.Debugf(template, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.SugaredLogger.Infof(template, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.SugaredLogger.Warnf(template, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.SugaredLogger.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.SugaredLogger.Errorf(template, args...)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.SugaredLogger.DPanic(args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.SugaredLogger.DPanicf(template, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.SugaredLogger.Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.SugaredLogger.Panicf(template, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.SugaredLogger.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.SugaredLogger.Fatalf(template, args...)
}

func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}
