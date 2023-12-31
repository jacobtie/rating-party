package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/handlers"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/logger"

	"github.com/rs/zerolog/log"
)

var (
	BUILD_TAG      = "unknown"
	BUILD_DATE     = "unknown"
	BUILD_GIT_HASH = "unknown"
)

//go:embed dist
var clientDir embed.FS

func main() {
	if err := run(); err != nil {
		log.Err(err).Msg("failed to run server")
	}
}

func run() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	logger.InitLogger(cfg.Environment, cfg.Instance, BUILD_TAG, BUILD_GIT_HASH, BUILD_DATE)
	logger.Get().Info().Msg("starting rating party server")
	db, err := db.New(cfg.DB.DBUser, cfg.DB.DBPass, cfg.DB.DBURI, cfg.DB.DBName)
	if err != nil {
		return err
	}
	strippedClientDir, err := fs.Sub(clientDir, "dist")
	if err != nil {
		return fmt.Errorf("failed to strip client directory prefix: %w", err)
	}
	server := &http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      handlers.NewAPI(cfg, db, strippedClientDir),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}
	logger.Get().Info().
		Str("environment", string(cfg.Environment)).
		Str("host", cfg.Web.APIHost).
		Msg("service started")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed while running server: %w", err)
	}
	logger.Get().Info().Msg("service stopped")
	return nil
}
