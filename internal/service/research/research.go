package research

import (
	"apiservice/internal/config"
	"io"
	"log/slog"
)

type InferenceProvider interface {
	DoInference(filepath string) error
}

type ResearchProvider interface {
	SaveFile(filename string, src io.Reader) error
}

type ResearchService struct {
	log               *slog.Logger
	cfg               *config.Config
	researchProvider  ResearchProvider
	inferenceProvider InferenceProvider
}

func New(log *slog.Logger, cfg *config.Config, researchProvider ResearchProvider, inferenceProvider InferenceProvider) *ResearchService {
	return &ResearchService{
		log:               log,
		cfg:               cfg,
		researchProvider:  researchProvider,
		inferenceProvider: inferenceProvider,
	}
}

func (s *ResearchService) SaveFile(filename string, src io.Reader) error {
	err := s.researchProvider.SaveFile(filename, src)
	if err != nil {
		s.log.Error("fail save file", slog.String("err", err.Error()))
		return err
	}

	s.log.Info("start inference")
	err = s.inferenceProvider.DoInference(s.cfg.ResearchSavePath + "/")
	if err != nil {
		s.log.Error("inference error", slog.String("err", err.Error()))
	}
	s.log.Info("finish inference")

	return nil
}
