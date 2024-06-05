package app

import (
	"fmt"

	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/Nol1feee/birthday-notifier/config"
	"github.com/Nol1feee/birthday-notifier/pkg/database/postgres"
	"github.com/Nol1feee/birthday-notifier/pkg/logger"
)

func Run(cfg *config.Config) {
	fmt.Printf("%+v\n", cfg)

	db, err := postgres.NewPostgresConnection(cfg.DB)
	if err != nil {
		logger.Fatal("error connecting to database", zap.Error(err))
	}

	//миграции

	defer db.Close()

	postgres.MigrateDB(db, cfg.DB)

	logger.Info("PG connected")
}
