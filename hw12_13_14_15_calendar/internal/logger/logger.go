package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logger struct { // TODO
	level    map[string]bool
	instance *zap.Logger
}

func New(level string, file string) *logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   file,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.TimeKey = "time"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		w,
		zap.InfoLevel,
	)

	levels := strings.Split(level, "|")
	levelMap := make(map[string]bool)

	for _, level := range levels {
		levelMap[level] = true
	}

	return &logger{levelMap, zap.New(core)}
}

func (l logger) Info(msg string, keysAndValues ...interface{}) {
	if l.level["info"] {
		l.instance.Sugar().Infow(msg, keysAndValues...)
	}
}

func (l logger) Debug(msg string, keysAndValues ...interface{}) {
	if l.level["debug"] {
		l.instance.Sugar().Infow(msg, keysAndValues...)
	}
}

func (l logger) Warn(msg string, keysAndValues ...interface{}) {
	if l.level["warn"] {
		l.instance.Sugar().Infow(msg, keysAndValues...)
	}
}

func (l logger) Error(msg string, keysAndValues ...interface{}) {
	if l.level["error"] {
		l.instance.Sugar().Infow(msg, keysAndValues...)
	}
}
