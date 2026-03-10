package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once

func init() {
	once.Do(func() {
		var logger *zap.Logger

		consoleDebugging := zapcore.Lock(os.Stdout)

		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "/go/logs/app.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
		})

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
		fileEncoder := zapcore.NewJSONEncoder(productionCfg)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel),
			zapcore.NewCore(fileEncoder, file, zapcore.InfoLevel),
		)

		logger = zap.New(core)
		defer logger.Sync()

		zap.ReplaceGlobals(logger)
	})
}
