package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"github.com/Nol1feee/birthday-notifier/pkg/database/postgres"
)

type Config struct {
	DB    postgres.DB
	HTTP  HTTP `yaml:"http"`
	Email `yaml:"email"`
}

type (
	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port int    `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Email struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		From string `yaml:"from"`
		Pass string `env-required:"true" env:"EMAIL_PASS"`
	}
)

var (
	cfg  = &Config{}
	once sync.Once
)

func New() *Config {
	once.Do(func() {
		// Загружаем переменные окружения из .env файла
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}

		// Загружаем конфигурацию из YAML файла
		if err := cleanenv.ReadConfig("./config/config.yaml", cfg); err != nil {
			log.Fatal("Error reading config.yaml file:", err)
		}

		// Перезаписываем значения конфигурации переменными окружения
		if err := cleanenv.ReadEnv(cfg); err != nil {
			log.Fatal("Error reading environment variables:", err)
		}
	})

	return cfg
}
