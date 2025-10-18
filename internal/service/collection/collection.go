package collection

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/collection"
	"apiservice/internal/domain/research"
	"log/slog"
	"path/filepath"

	"github.com/google/uuid"
)

type collectionProvider interface {
	Create(id string, pathologyLevel float32) (collection.Collection, error)
	Delete(id string) error
	WritePathologyLevel(id string, pathologyLevel float32) error
	WriteTitle(id string, title string) error
	List() ([]collection.Collection, error)
	GetOne(id string) (collection.Collection, error)
	CheckExists(id string) (bool, error)
}

type researchProvider interface {
	List(collectionId string) ([]research.Research, error)
	DeleteFiles(collectionId string) error
}

type Service struct {
	log                *slog.Logger
	cfg                *config.Config
	collectionProvider collectionProvider
	researchProvider   researchProvider
}

func New(log *slog.Logger, cfg *config.Config, collectionProvider collectionProvider, researchProvider researchProvider) *Service {
	return &Service{
		log:                log,
		cfg:                cfg,
		collectionProvider: collectionProvider,
		researchProvider:   researchProvider,
	}
}

func (s *Service) Create() (collection.Collection, error) {
	id := uuid.New().String()[:13] // чтобы имена папок с коллекциями были достаточно короткими
	res, err := s.collectionProvider.Create(id, s.cfg.PathologyLevel)
	if err != nil {
		s.log.Error("create collection", slog.String("err", err.Error()))
		return collection.Collection{}, err
	}

	return res, nil
}

func (s *Service) Delete(id string) error {
	if err := s.researchProvider.DeleteFiles(id); err != nil {
		s.log.Error("delete collection files", slog.String("id", id), slog.String("err", err.Error()))
		return err
	}

	if err := s.collectionProvider.Delete(id); err != nil {
		s.log.Error("delete collection", slog.String("id", id), slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *Service) Update(id string, update collection.Update) error {
	if update.PathologyLevel != nil {
		s.log.Info("update pathologyLevel of collection")
		if err := s.collectionProvider.WritePathologyLevel(id, *update.PathologyLevel); err != nil {
			s.log.Error("fail set pathology level", slog.String("collectionId", id))
			return err
		}
	}

	if update.Title != nil {
		s.log.Info("update title of collection")
		if err := s.collectionProvider.WriteTitle(id, *update.Title); err != nil {
			s.log.Error("fail set title level", slog.String("collectionId", id))
			return err
		}
	}

	return nil
}

func (s *Service) List() ([]collection.Collection, error) {
	list, err := s.collectionProvider.List()
	if err != nil {
		s.log.Error("list collection", slog.String("err", err.Error()))
		return nil, err
	}

	return list, err
}

func (s *Service) GetOne(id string) (collection.CollectionWithResearches, error) {
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
			Size:                   rs[k].Size,
			Assessment:             rs[k].Assessment,
			ArchiveCorrupt:         rs[k].ArchiveCorrupt,
			ProbabilityOfPathology: rs[k].ProbabilityOfPathology,
			CreatedAt:              rs[k].CreatedAt,
			ProcessingStartedAt:    rs[k].ProcessingStartedAt,
			ProcessingFinishedAt:   rs[k].ProcessingFinishedAt,
			InferenceError:         rs[k].InferenceError,
			Metadata:               rs[k].Metadata,
			Filename:               filename,
			ProcessingDuration:     int64(diff.Seconds()),
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

func (s *Service) CheckExists(id string) (bool, error) {
	exists, err := s.collectionProvider.CheckExists(id)
	if err != nil {
		s.log.Error("check exists", slog.String("id", id), slog.String("err", err.Error()))
		return false, err
	}

	return exists, nil
}
