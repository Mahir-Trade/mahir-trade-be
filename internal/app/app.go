package app

import (
	"context"
	"database/sql"
	"mahir-trade-be/internal/app/infra"
	"mahir-trade-be/pkg/di"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var exitSigs = []os.Signal{syscall.SIGTERM, syscall.SIGINT}

func Start() {
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, exitSigs...)

	go func() {
		defer func() { exitCh <- syscall.SIGTERM }()
		if err := di.Invoke(startApp); err != nil {
			log.Fatal().Msgf("startApp: %s", err.Error())
		}
	}()
	<-exitCh

	if err := di.Invoke(gracefulShutdown); err != nil {
		log.Error().Msgf("gracefulShutdown: %s", err.Error())
	}
}

func startApp(
	e *echo.Echo,
	appCfg *infra.AppCfg,
) error {
	if err := di.Invoke(setRoute); err != nil {
		return err
	}

	return e.StartServer(&http.Server{
		Addr:         appCfg.Address,
		ReadTimeout:  appCfg.ReadTimeout,
		WriteTimeout: appCfg.WriteTimeout,
	})
}

func gracefulShutdown(
	e *echo.Echo,
	pg *sql.DB,
) {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Info().Msg("shutting down server")

	if err := pg.Close(); err != nil {
		log.Error().Msgf("postgres close: %s", err.Error())
	}

	if err := e.Shutdown(ctx); err != nil {
		log.Error().Msgf("echo shutdown: %s", err.Error())
	}

	log.Info().Msg("server gracefully stopped")
}