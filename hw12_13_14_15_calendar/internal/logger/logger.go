package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
	file   *os.File
}

func GetLogger(level string) (*ZapLogger, error) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "panic":
		zapLevel = zapcore.PanicLevel
	case "fatal":
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
	return context.WithValue(ctx, "logger", logger)
}

func GetLoggerFromContext(ctx context.Context) (*ZapLogger, error) {
	if l, ok := ctx.Value("logger").(*ZapLogger); ok {
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
	l.logger.Info(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.Warn(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.Error(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Fatal(msg string, fields map[string]interface{}) {
	l.logger.Fatal(msg, zap.Any("args", fields))
}
