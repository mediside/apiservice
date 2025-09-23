package service

import (
	"apiservice/internal/config"
	"apiservice/internal/service/collection"
	"apiservice/internal/service/research"
	"apiservice/internal/storage"
	"log/slog"
)

type Service struct {
	CollectionService *collection.CollectionService
	ResearchService   *research.ResearchService
}

func New(log *slog.Logger, cfg *config.Config, repo storage.Storage) Service {

	col := collection.New(log, cfg, repo.CollectionStorage)
	res := research.New(log, cfg, repo.ResearchStorage)

	return Service{
		CollectionService: col,
		ResearchService:   res,
	}
}
