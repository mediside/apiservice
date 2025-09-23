package collection

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/collection"
	"log/slog"

	"github.com/google/uuid"
)

type CollectionProvider interface {
	Create(id string, pathologyLevel float32) (collection.Collection, error)
	Delete(id string) error
	List() ([]collection.Collection, error)
	GetOne(id string) (collection.Collection, error)
}

type CollectionService struct {
	log                *slog.Logger
	cfg                *config.Config
	collectionProvider CollectionProvider
}

func New(log *slog.Logger, cfg *config.Config, collectionProvider CollectionProvider) *CollectionService {
	return &CollectionService{
		log:                log,
		cfg:                cfg,
		collectionProvider: collectionProvider,
	}
}

func (s *CollectionService) Create() (collection.Collection, error) {
	id := uuid.New()
	res, err := s.collectionProvider.Create(id.String(), s.cfg.PathologyLevel)
	if err != nil {
		s.log.Error("create collection", slog.String("err", err.Error()))
		return collection.Collection{}, err
	}

	return res, nil
}

func (s *CollectionService) Delete(id string) error {
	if err := s.collectionProvider.Delete(id); err != nil {
		s.log.Error("delete collection", slog.String("err", err.Error()))
	}

	return nil
}

func (s *CollectionService) List() ([]collection.Collection, error) {
	list, err := s.collectionProvider.List()
	if err != nil {
		s.log.Error("list collection", slog.String("err", err.Error()))
		return nil, err
	}

	return list, err
}

func (s *CollectionService) GetOne(id string) (collection.Collection, error) {
	res, err := s.collectionProvider.GetOne(id)
	if err != nil {
		s.log.Error("list collection", slog.String("err", err.Error()))
		return collection.Collection{}, err
	}

	return res, err
}
