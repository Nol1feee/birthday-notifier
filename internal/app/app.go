package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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

	/*INIT pg db */
	db, err := postgres.NewPostgresConnection(cfg.DB)
	if err != nil {
		logger.Fatal("error connecting to database", zap.Error(err))
	}

	defer db.Close()
	logger.Info("DB connected")

	postgres.MigrateDB(db, cfg.DB)

	/*INIT services*/
	usersRepo := psql.NewUsers(db)
	usersService := service.NewUsers(usersRepo)
	notifierService := service.NewNotifier(cfg.Email, usersRepo)
	handler := rest.NewHandler(usersService)

	/*INIT http server */
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           handler.InitRouter(),
		ReadHeaderTimeout: 0,
	}

	logger.Info(fmt.Sprintf(srv.Addr))

	/*INIT cron notifier */
	c := cron.New(cron.WithLocation(time.FixedZone("MSK", 3*60*60)))

	_, err = c.AddFunc("0 7 * * *", func() {
		err := notifierService.CongratulateAll()
		if err != nil {
			logger.Error(err.Error())
		}
	})

	if err != nil {
		logger.Error("Ошибка при отправке email'ов", zap.Error(err))
	}

	go func() {
		c.Start()
		logger.Info("cron service started")
		defer c.Stop()

		select {}
	}()

	/*RUN http server */
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit

	//нормально обработать grace -> srv.shutdown и тд
	logger.Info("signal received to end the program")
}
