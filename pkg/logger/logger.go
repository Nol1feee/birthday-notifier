package logger

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	Mode string `yaml:"log_mode" env-default:"dev"`
}

const (
	envDev  = "dev"
	envProd = "prod"
)

var globalLogger *zap.Logger

func init() {
	log := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", log)
	if err != nil {
		fmt.Printf("logger error - %s\n", err)
		os.Exit(1)
	}

	var cfg zap.Config

	switch log.Mode {
	case envDev:
		cfg = zap.NewDevelopmentConfig()
	case envProd:
		cfg = zap.NewProductionConfig()
	default:
		fmt.Printf("unknown logger modde, expected '%s' or '%s'", envDev, envProd)
		return
	}

	cfg.DisableStacktrace = true
	globalLogger, _ = cfg.Build(zap.AddCallerSkip(1))
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
