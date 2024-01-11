package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
	file   *os.File
}

type KeyLogger string

func GetLogger(level string) (*ZapLogger, error) {
	level = strings.ToUpper(level)
	var zapLevel zapcore.Level
	switch level {
	case "DEBUG":
		zapLevel = zapcore.DebugLevel
	case "INFO":
		zapLevel = zapcore.InfoLevel
	case "WARN":
		zapLevel = zapcore.WarnLevel
	case "ERROR":
		zapLevel = zapcore.ErrorLevel
	case "PANIC":
		zapLevel = zapcore.PanicLevel
	case "FATAL":
		zapLevel = zapcore.FatalLevel
	default:
		return nil, fmt.Errorf("unsupported level of logger: %s", level)
	}

	err := os.MkdirAll("logs", 0o777)
	if err != nil {
		return nil, fmt.Errorf("could not create directory %w", err)
	}

	file, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o777)
	if err != nil {
		return nil, fmt.Errorf("could not open file %w", err)
	}

	logger := zap.New(zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(file),
			zapLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zapLevel,
		),
	),
	)

	logg := logger.Sugar()

	return &ZapLogger{logg, file}, nil
}

func ContextWithLogger(ctx context.Context, logger *ZapLogger) context.Context {
	return context.WithValue(ctx, KeyLogger("logger"), logger)
}

func GetLoggerFromContext(ctx context.Context) (*ZapLogger, error) {
	if l, ok := ctx.Value(KeyLogger("logger")).(*ZapLogger); ok {
		return l, nil
	}

	return GetLogger("info")
}

func (l *ZapLogger) Close() {
	err := l.file.Close()
	if err != nil {
		log.Fatalf("error while closing file with logs: %v", err)
	}
}

func (l *ZapLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.Debug(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.Infow(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.Warnw(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.Errorw(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Fatal(msg string, fields map[string]interface{}) {
	l.logger.Fatalw(msg, zap.Any("args", fields))
}
