package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ocmaster/backend/application"
	mysqlrepo "github.com/ocmaster/backend/infrastructure/persistence/mysql"
	httpx "github.com/ocmaster/backend/interfaces/http"
	"github.com/ocmaster/backend/interfaces/http/handler"
	"github.com/ocmaster/backend/pkg/config"
)

var Version = "0.0.1" // 注入: go build -ldflags="-X main.Version=$(cat VERSION)"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Str("version", Version).Msg("starting")

	cfg := config.Load()

	db, err := cfg.OpenDB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database")
	}
	if err := mysqlrepo.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	// Repositories
	hardwareRepo := mysqlrepo.NewHardwareRepo(db)
	sharecodeRepo := mysqlrepo.NewShareCodeRepo(db)
	merchantRepo := mysqlrepo.NewMerchantRepo(db)
	suggestionRepo := mysqlrepo.NewSuggestionRepo(db)
	referenceRepo := mysqlrepo.NewReferenceRepo(db)

	// Services
	hardwareSvc := application.NewHardwareService(hardwareRepo, sharecodeRepo)
	merchantSvc := application.NewMerchantService(merchantRepo, cfg.JWTSecret)
	suggestionSvc := application.NewSuggestionService(suggestionRepo)
	referenceSvc := application.NewReferenceService(referenceRepo)
	cleanupSvc := application.NewCleanupService(hardwareRepo, suggestionRepo)

	// Handlers
	hardwareH := handler.NewHardwareHandler(hardwareSvc)
	merchantH := handler.NewMerchantHandler(merchantSvc)
	suggestionH := handler.NewSuggestionHandler(suggestionSvc)
	referenceH := handler.NewReferenceHandler(referenceSvc)

	router := httpx.NewRouter(httpx.RouterDeps{
		HardwareHandler:   hardwareH,
		MerchantHandler:   merchantH,
		SuggestionHandler: suggestionH,
		ReferenceHandler:  referenceH,
		CleanupService:    cleanupSvc,
		JWTSecret:         cfg.JWTSecret,
	})

	// Start cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cleanupSvc.Start(ctx)

	srv := &http.Server{Addr: ":" + cfg.ServerPort, Handler: router}

	go func() {
		log.Info().Str("port", cfg.ServerPort).Msg("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}
	log.Info().Msg("server stopped")
}
