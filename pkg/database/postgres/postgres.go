package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Nol1feee/birthday-notifier/pkg/logger"
)

const (
	pathMigrations = "file://migrations"
)

type DB struct {
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Name     string `yaml:"name"  env:"DB_NAME" env-default:"postgres"`
	User     string `yaml:"user"  env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
	SSLMode  string `yaml:"sslmode"  env:"DB_SSLMODE" env-default:"disable"`
}

// NewPostgresConnection creates and verifies a PostgreSQL database connection based on cfg.
// On success it returns an opened *sql.DB ready for use. If the database cannot be opened
// or connectivity cannot be verified via ping, it returns a non-nil error describing the failure.
func NewPostgresConnection(cfg DB) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.SSLMode, cfg.Password)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("database open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	return db, nil
}

// MigrateDB runs file-based migrations from pathMigrations against the database named in cfg using the provided *sql.DB.
// It logs a fatal error and terminates the process if driver creation, migration setup, or the migration run fails for reasons other than migrate.ErrNoChange; if migrations are applied successfully it logs an informational message.
func MigrateDB(db *sql.DB, cfg DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(pathMigrations, cfg.Name, driver)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("", zap.Error(err))
	} else if err == nil {
		logger.Info("Database migration was run successfully")
	}
}
