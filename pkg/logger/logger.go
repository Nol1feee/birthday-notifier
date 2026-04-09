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

var globalLogger *zap.Logger

// init initializes the package logger by reading configuration from ./config/config.yaml
// and configuring the global logger according to the Config.Mode.
// It prints an error and exits the process if configuration cannot be read; if the mode is unrecognized
// it prints an error and leaves globalLogger unset.
func init() {
	const (
		envDev      = "dev"
		envProd     = "prod"
		defaultPath = "./config/config.yaml"
	)

	log := &Config{}

	err := cleanenv.ReadConfig(defaultPath, log)
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
