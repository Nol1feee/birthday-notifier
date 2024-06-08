package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/Nol1feee/birthday-notifier/config"
	"github.com/Nol1feee/birthday-notifier/internal/service"
	"github.com/Nol1feee/birthday-notifier/internal/storage/psql"
	"github.com/Nol1feee/birthday-notifier/internal/transport/rest"
	"github.com/Nol1feee/birthday-notifier/pkg/database/postgres"
	"github.com/Nol1feee/birthday-notifier/pkg/logger"
)

func Run(cfg *config.Config) {
	logger.Debug("cfg info", zap.String("cfg", fmt.Sprintf("%+v", cfg)))

	/* INIT pg db */
	db, err := postgres.NewPostgresConnection(cfg.DB)
	if err != nil {
		logger.Fatal("error connecting to database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("DB connected")

	postgres.MigrateDB(db, cfg.DB)

	/* INIT services */
	usersRepo := psql.NewUsers(db)
	usersService := service.NewUsers(usersRepo)
	notifierService := service.NewNotifier(cfg.Email, usersRepo)
	notifierService.NotifyingUpcomingBirthdays()

	handler := rest.NewHandler(usersService)

	/* INIT http server */
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           handler.InitRouter(),
		ReadHeaderTimeout: 0,
	}

	logger.Info(fmt.Sprintf("Starting HTTP server at %s", srv.Addr))

	/* INIT cron notifier */
	c := cron.New(cron.WithLocation(time.FixedZone("MSK", 3*60*60)))
	_, err = c.AddFunc("0 7 * * *", func() {
		if err := notifierService.CongratulateAll(); err != nil {
			logger.Error("Error in CongratulateAll", zap.Error(err))
		}
		if err := notifierService.NotifyingUpcomingBirthdays(); err != nil {
			logger.Error("Error in NotifyingUpcomingBirthdays", zap.Error(err))
		}
	})

	if err != nil {
		logger.Error("Error scheduling cron job", zap.Error(err))
		return
	}

	c.Start()
	defer c.Stop()

	go func() {
		logger.Info("cron service started")
		select {}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit

	logger.Info("Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("HTTP server forced to shutdown", zap.Error(err))
	}

	logger.Info("HTTP server stopped successfully")
}
