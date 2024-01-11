package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"C:/Users/1/Desktop/test.log"}
	cfg.ErrorOutputPaths = []string{"C:/Users/1/Desktop/test.log"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic("zap logger init fail!")
	}
}
