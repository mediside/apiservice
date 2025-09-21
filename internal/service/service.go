package service

import (
	"apiservice/internal/config"
	"apiservice/internal/service/research"
	"apiservice/internal/storage"
	"log/slog"
)

type Service struct {
	ResearchService *research.ResearchService
}

func New(log *slog.Logger, cfg *config.Config, repo storage.Storage) Service {

	res := research.New(log, cfg, repo.ResearchStorage)

	return Service{
		ResearchService: res,
	}
}
