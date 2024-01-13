package utils

import (
	"MyOSS/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(config.LOG_LEVEL)
	cfg.OutputPaths = []string{config.LOG_DIR}
	cfg.ErrorOutputPaths = []string{config.LOG_DIR}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic("zap logger init fail!")
	}
}
