package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func init() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/grinnodes.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     60,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		zap.InfoLevel,
	)
	Logger = zap.New(core)
}
