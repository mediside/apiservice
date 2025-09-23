package research

import (
	"apiservice/internal/config"
	"io"
	"log/slog"
)

type ResearchProvider interface {
	SaveFile(filename string, src io.Reader) error
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

func (s *ResearchService) SaveFile(filename string, src io.Reader) error {
	err := s.researchProvider.SaveFile(filename, src)
	if err != nil {
		s.log.Error("fail save file", slog.String("err", err.Error()))
		return err
	}

	return nil
}
