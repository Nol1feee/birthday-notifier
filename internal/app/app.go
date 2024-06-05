package app

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/Nol1feee/birthday-notifier/config"
	"github.com/Nol1feee/birthday-notifier/internal/service"
	"github.com/Nol1feee/birthday-notifier/internal/storage/psql"
	"github.com/Nol1feee/birthday-notifier/internal/transport/rest"
	"github.com/Nol1feee/birthday-notifier/pkg/database/postgres"
	"github.com/Nol1feee/birthday-notifier/pkg/logger"
)

func Run(cfg *config.Config) {
	fmt.Printf("%+v\n", cfg)

	db, err := postgres.NewPostgresConnection(cfg.DB)
	if err != nil {
		logger.Fatal("error connecting to database", zap.Error(err))
	}

	defer db.Close()
	logger.Info("DB connected")

	postgres.MigrateDB(db, cfg.DB)

	usersRepo := psql.NewUsers(db)
	usersService := service.NewUsers(usersRepo)
	handler := rest.NewHandler(usersService)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           handler.InitRouter(),
		ReadHeaderTimeout: 0,
	}

	logger.Info(fmt.Sprintf(srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}

	//gracefull shutdown
}
