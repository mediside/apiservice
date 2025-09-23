package service

import (
	"apiservice/internal/config"
	"apiservice/internal/service/collection"
	"apiservice/internal/storage"
	"log/slog"
)

type Service struct {
	CollectionService *collection.CollectionService
}

func New(log *slog.Logger, cfg *config.Config, repo storage.Storage) Service {

	res := collection.New(log, cfg, repo.CollectionStorage)

	return Service{
		CollectionService: res,
	}
}
