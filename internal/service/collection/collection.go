package collection

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/collection"
	"apiservice/internal/domain/research"
	"log/slog"
	"path/filepath"

	"github.com/google/uuid"
)

type CollectionProvider interface {
	Create(id string, pathologyLevel float32) (collection.Collection, error)
	Delete(id string) error
	List() ([]collection.Collection, error)
	GetOne(id string) (collection.Collection, error)
	CheckExists(id string) (bool, error)
}

type ResearchProvider interface {
	List(collectionId string) ([]research.Research, error)
}

type CollectionService struct {
	log                *slog.Logger
	cfg                *config.Config
	collectionProvider CollectionProvider
	researchProvider   ResearchProvider
}

func New(log *slog.Logger, cfg *config.Config, collectionProvider CollectionProvider, researchProvider ResearchProvider) *CollectionService {
	return &CollectionService{
		log:                log,
		cfg:                cfg,
		collectionProvider: collectionProvider,
		researchProvider:   researchProvider,
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
		return err
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

func (s *CollectionService) GetOne(id string) (collection.CollectionWithResearches, error) {
	col, err := s.collectionProvider.GetOne(id)
	if err != nil {
		s.log.Error("get one collection", slog.String("err", err.Error()))
		return collection.CollectionWithResearches{}, err
	}

	rs, err := s.researchProvider.List(id)
	if err != nil {
		s.log.Error("list researches", slog.String("err", err.Error()))
		return collection.CollectionWithResearches{}, err
	}

	researches := make([]research.ResearchResult, len(rs))

	for k := range rs {
		filename := filepath.Base(rs[k].Filepath)
		diff := rs[k].ProcessingFinishedAt.Sub(rs[k].ProcessingStartedAt)

		researches[k] = research.ResearchResult{
			Id:                     rs[k].Id,
			Filepath:               rs[k].Filepath,
			Assessment:             rs[k].Assessment,
			ArchiveCorrupt:         rs[k].ArchiveCorrupt,
			ProbabilityOfPathology: rs[k].ProbabilityOfPathology,
			CreatedAt:              rs[k].CreatedAt,
			ProcessingStartedAt:    rs[k].ProcessingStartedAt,
			ProcessingFinishedAt:   rs[k].ProcessingFinishedAt,
			InferenceError:         rs[k].InferenceError,
			Metadata:               rs[k].Metadata,
			Filename:               filename,
			ProcessingDuration:     diff.Milliseconds(),
		}
	}

	fullColl := collection.CollectionWithResearches{
		Id:             col.Id,
		Num:            col.Num,
		Title:          col.Title,
		PathologyLevel: col.PathologyLevel,
		CreatedAt:      col.CreatedAt,
		Researches:     researches,
	}

	return fullColl, err
}

func (s *CollectionService) CheckExists(id string) (bool, error) {
	exists, err := s.collectionProvider.CheckExists(id)
	if err != nil {
		s.log.Error("check exists", slog.String("id", id), slog.String("err", err.Error()))
		return false, err
	}

	return exists, nil
}
