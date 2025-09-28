package research

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/inference"
	"apiservice/internal/domain/research"
	"io"
	"log/slog"
	"time"
)

type InferenceProvider interface {
	DoInference(responseCh chan<- inference.InferenceResponse, filepath, studyId, seriesId string) error
}

type ResearchProvider interface {
	SaveFile(collectionId, filename string, src io.Reader) error
	Create(id, collectionId, filepath string, size int64, archiveCorrupt bool, metadata research.ResearchMetadata) error
	Delete(id string) error
	WriteInferenceResult(id string, probabilityOfPathology float32) error
	WriteInferenceError(id, inferenceErr string) error
	WriteInferenceFinishTime(id string, finishedAt time.Time) error
	WriteInferenceStartTime(id string, startedAt time.Time) error
}

type ResearchService struct {
	log               *slog.Logger
	cfg               *config.Config
	researchProvider  ResearchProvider
	inferenceProvider InferenceProvider
	taskCh            chan inference.InferenceTask
	inferenceCh       chan inference.InferenceProgress // для отправки во внешний мир
}

func New(log *slog.Logger, cfg *config.Config, researchProvider ResearchProvider, inferenceProvider InferenceProvider) *ResearchService {
	research := &ResearchService{
		log:               log,
		cfg:               cfg,
		researchProvider:  researchProvider,
		inferenceProvider: inferenceProvider,
		taskCh:            make(chan inference.InferenceTask),
		inferenceCh:       make(chan inference.InferenceProgress),
	}

	go research.inferenceWorker()

	return research
}

func (s *ResearchService) RunFileProcessing(filename, collectionId string, src io.Reader) error {
	// загруженный файл сохраняем в любом случае, чтобы он отображался в статистике
	// даже если он не валидный, чтобы было видно пользователю, что файл загрузился, но не читается
	err := s.researchProvider.SaveFile(collectionId, filename, src)
	if err != nil {
		s.log.Error("fail save file", slog.String("filename", filename), slog.String("err", err.Error()))
		return err
	}

	// горутина нужна, чтобы не блокировать HTTP-вызов
	// она сохранит контекст и встанет дожидаться очереди на отправку задачи
	go s.processing(filename, collectionId)

	return nil
}

func (s *ResearchService) Delete(id string) error {
	if err := s.researchProvider.Delete(id); err != nil {
		s.log.Error("delete collection", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *ResearchService) InferenceCh() <-chan inference.InferenceProgress {
	return s.inferenceCh
}
