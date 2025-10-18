package research

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/inference"
	"apiservice/internal/domain/research"
	"io"
	"log/slog"
	"time"
)

type inferenceProvider interface {
	DoInference(responseCh chan<- inference.InferenceResponse, filepath, studyId, seriesId string) error
}

type researchProvider interface {
	SaveFile(collectionId, filename string, src io.Reader) error
	Create(id, collectionId, filepath string, size int64, archiveCorrupt bool, metadata research.Metadata) error
	DeleteEntry(filepath string) error
	DeleteSingleFile(filepath string) error
	GetFilepath(id string) (string, error)
	CheckExists(collectionId, filename string) (bool, error)
	WriteInferenceResult(id string, probabilityOfPathology float32) error
	WriteInferenceError(id, inferenceErr string) error
	WriteInferenceFinishTime(id string, finishedAt time.Time) error
	WriteInferenceStartTime(id string, startedAt time.Time) error
}

type Service struct {
	log               *slog.Logger
	cfg               *config.Config
	researchProvider  researchProvider
	inferenceProvider inferenceProvider
	taskCh            chan inference.InferenceTask
	inferenceCh       chan inference.InferenceProgress // для отправки во внешний мир
	updateCh          chan research.ResearchUpdate     // для общих обновлений в БД (кроме удаления)
}

func New(log *slog.Logger, cfg *config.Config, researchProvider researchProvider, inferenceProvider inferenceProvider) *Service {
	research := &Service{
		log:               log,
		cfg:               cfg,
		researchProvider:  researchProvider,
		inferenceProvider: inferenceProvider,
		taskCh:            make(chan inference.InferenceTask),
		inferenceCh:       make(chan inference.InferenceProgress),
		updateCh:          make(chan research.ResearchUpdate),
	}

	go research.inferenceWorker()

	return research
}

func (s *Service) RunFileProcessing(filename, collectionId string, src io.Reader) error {
	// загруженный файл сохраняем в любом случае, чтобы он отображался в статистике
	// даже если он не валидный, чтобы было видно пользователю, что файл загрузился, но не читается
	if err := s.researchProvider.SaveFile(collectionId, filename, src); err == research.ErrFileAlreadyExists {
		s.log.Warn("file alrady exists; skip it", slog.String("filename", filename), slog.String("collectionId", collectionId))
		return nil
	} else if err != nil {
		s.log.Error("fail save file", slog.String("filename", filename), slog.String("err", err.Error()))
		return err
	}

	s.log.Info("success upload file", slog.String("filename", filename), slog.String("collectionId", collectionId))

	// горутина нужна, чтобы не блокировать HTTP-вызов
	// она сохранит контекст и встанет дожидаться очереди на отправку задачи
	go s.processing(filename, collectionId)

	return nil
}

func (s *Service) Delete(id string) error {
	filepath, err := s.researchProvider.GetFilepath(id)
	if err != nil {
		s.log.Error("fail get research filepath", slog.String("id", id), slog.String("err", err.Error()))
		return err
	}

	if err := s.researchProvider.DeleteEntry(filepath); err != nil {
		s.log.Error("fail delete research entry", slog.String("id", id), slog.String("err", err.Error()))
		return err
	}

	return nil // удаление файла произойдет автоматически после инференса
}

func (s *Service) CheckExists(collectionId, filename string) (bool, error) {
	exists, err := s.researchProvider.CheckExists(collectionId, filename)
	if err != nil {
		s.log.Error("check exists", slog.String("collectionId", collectionId), slog.String("filename", filename), slog.String("err", err.Error()))
		return false, err
	}

	return exists, nil
}

func (s *Service) InferenceCh() <-chan inference.InferenceProgress {
	return s.inferenceCh
}

func (s *Service) UpdateCh() <-chan research.ResearchUpdate {
	return s.updateCh
}
