package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joeshaw/envdecode"
	"github.com/oklog/run"
	"github.com/rs/zerolog/log"

	"github.com/dagowa/adservice/internal/controllers"
	"github.com/dagowa/adservice/internal/server"
	"github.com/dagowa/adservice/internal/storage"
	"github.com/dagowa/adservice/pkg/logger"
)

type config struct {
	Logger  logger.Config
	Service server.Config
}

func main() {
	cfg := &config{}
	if err := envdecode.StrictDecode(cfg); err != nil {
		log.Fatal().Err(err).Msg("Cannot decode config envs")
	}

	l := logger.NewLogger(&cfg.Logger)
	ctx := l.WithContext(context.Background())
	l.Info().Interface("config", cfg).Msg("The gathered config")

	psqlConn, err := storage.New().NewPostgreSQLConn(cfg.Service.PSQLConfig)
	if err != nil {
		l.Fatal().Err(err).Msg("Cannot set up psql connection")
	}
	defer psqlConn.Close()

	controllerSet := controllers.NewControllerSet(psqlConn.Pool())

	ctx, cancel := context.WithCancel(l.WithContext(ctx))

	g := &run.Group{}
	{
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error {
			<-stop
			return nil
		}, func(error) {
			signal.Stop(stop)
			cancel()
			close(stop)
		})
	}
	{
		srv := server.NewServer(ctx, &cfg.Service, &controllerSet)

		g.Add(func() error {
			l.Info().Str("address", srv.Addr).Msg("Start listening")
			if err := srv.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					return nil
				}
				return err
			}
			l.Info().Msg("Listening is stopped")
			return nil
		}, func(error) {
			l.Info().Msg("Shutdowning listening...")
			if err := srv.Shutdown(ctx); err != nil {
				l.Error().Err(err).Msg("Cannot shutdown the service properly")
			}
			l.Info().Msg("The service is shutdown")
		})
	}

	l.Info().Msg("Running the service...")
	if err := g.Run(); err != nil {
		l.Fatal().Err(err).Msg("The service has been stopped with error")
	}
	l.Info().Msg("The service is stopped")
}
