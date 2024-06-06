package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Entry *zap.Logger
)

func NewLogger() error {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		return err
	}

	logger = logger.WithOptions(zap.AddStacktrace(zap.DPanicLevel))

	Entry = logger

	return nil
}

func NewTestLogger() *zap.Logger {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	discardSyncer := zapcore.Lock(zapcore.AddSync(io.Discard))
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, discardSyncer, zap.InfoLevel),
	)

	return zap.New(core)
}
