package logger

import (
	"go.uber.org/zap"
	"log"
)

var (
	logger *zap.Logger
)

func init() {

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		Encoding:          "json",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		DisableStacktrace: true,
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		log.Fatal(err)
	}
}

func Error(err error, fields ...zap.Field) {
	logger.Error(err.Error(), fields...)
}

func Fatal(err error, fields ...zap.Field) {
	logger.Fatal(err.Error(), fields...)
}

func Sync() {
	err := logger.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
