package logger

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	envName   = "ENV"
	envDev    = "dev"
	envProd   = "prod"
	pathToEnv = "../../.env"
)

var globalLogger *zap.Logger

func init() {
	_ = godotenv.Load(pathToEnv)

	env := os.Getenv(envName)
	switch env {
	case envDev:
		globalLogger, _ = zap.NewDevelopment(zap.AddCallerSkip(1))
	case envProd:
		globalLogger, _ = zap.NewProduction(zap.AddCallerSkip(1))
	}

}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
	os.Exit(1)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}
