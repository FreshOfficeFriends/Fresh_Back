package logger

import (
	"fmt"

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
	fmt.Println("env = ", env)
	switch env {
	case envDev:
		globalLogger, _ = zap.NewDevelopment(zap.AddCallerSkip(1))
	case envProd:
		globalLogger, _ = zap.NewProduction(zap.AddCallerSkip(1))
	}
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg)
	os.Exit(1)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg)
}

//// without runtime.caller
//func InfoMiddleware(msg string, fields ...zap.Field) {
//	globalLogger.Info(msg, fields...)
//}

//func getCallerInfo() zap.Field {
//	_, file, line, _ := runtime.Caller(2)
//	return zap.String("caller", fmt.Sprintf("%s:%s", file, strconv.Itoa(line)))
//}
