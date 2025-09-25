package research

import (
	"apiservice/internal/config"
	"io"
	"log/slog"

	"github.com/google/uuid"
)

type InferenceProvider interface {
	DoInference(filepath string) error
}

type ResearchProvider interface {
	SaveFile(collectionId, filename string, src io.Reader) error
	Create(id, collectionId, filepath string) error
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

func (s *ResearchService) RunFileProcessing(filename, collectionId string, src io.Reader) error {
	err := s.researchProvider.SaveFile(collectionId, filename, src)
	if err != nil {
		s.log.Error("fail save file", slog.String("filename", filename), slog.String("err", err.Error()))
		return err
	}

	id := uuid.New()
	filepath := s.cfg.ResearchSavePath + "/" + filename
	err = s.researchProvider.Create(id.String(), collectionId, filepath)
	if err != nil {
		s.log.Error("fail create row in db", slog.String("err", err.Error()))
	}

	go s.Inference(filepath)

	return nil
}

func (s *ResearchService) Inference(filepath string) {
	s.log.Info("start inference", slog.String("filepath", filepath))

	err := s.inferenceProvider.DoInference(filepath)
	if err != nil {
		s.log.Error("inference error", slog.String("err", err.Error()))
	}
	s.log.Info("finish inference")
}
