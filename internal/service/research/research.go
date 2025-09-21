package research

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/research"
	"log/slog"

	"github.com/google/uuid"
)

type ResearchProvider interface {
	Create(id string, pathologyLevel float32) (research.Research, error)
}

type ResearchService struct {
	log              *slog.Logger
	cfg              *config.Config
	researchProvider ResearchProvider
}

func New(log *slog.Logger, cfg *config.Config, researchProvider ResearchProvider) *ResearchService {
	return &ResearchService{
		log:              log,
		cfg:              cfg,
		researchProvider: researchProvider,
	}
}

func (s *ResearchService) Create() (research.Research, error) {
	id := uuid.New()
	res, err := s.researchProvider.Create(id.String(), s.cfg.PathologyLevel)
	if err != nil {
		s.log.Error("create research", slog.String("err", err.Error()))
		return research.Research{}, err
	}

	return res, nil
}
