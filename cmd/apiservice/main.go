package main

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/model"
	"apiservice/internal/handler"
	"apiservice/internal/service"
	"apiservice/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting apiservice", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	repo := storage.New(log, cfg)
	defer repo.Close()

	serv := service.New(log, cfg, repo)

	hand := handler.New(log, cfg, serv)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	addr := fmt.Sprintf("%s:%s", cfg.Http.Host, cfg.Http.Port)
	srv := http.Server{Addr: addr, Handler: hand}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("server launch was interrupted", slog.String("err", err.Error()))
		}
	}()
	log.Info("server started", slog.String("addr", addr))

	<-stop
	log.Info("stopping server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.String("error", err.Error()))
	}
	log.Info("server stopped. Exit")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case model.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case model.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
