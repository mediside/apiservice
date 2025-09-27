package research

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/inference"
	"io"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type InferenceProvider interface {
	DoInference(responseCh chan<- inference.InferenceResponse, filepath string) error
}

type ResearchProvider interface {
	SaveFile(collectionId, filename string, src io.Reader) error
	Create(id, collectionId, filepath string) error
	WriteInferenceResult(id string, probabilityOfPathology float32, finishedAt time.Time) error
}

type ResearchService struct {
	log               *slog.Logger
	cfg               *config.Config
	researchProvider  ResearchProvider
	inferenceProvider InferenceProvider
	taskCh            chan inference.InferenceTask
}

func New(log *slog.Logger, cfg *config.Config, researchProvider ResearchProvider, inferenceProvider InferenceProvider) *ResearchService {
	research := &ResearchService{
		log:               log,
		cfg:               cfg,
		researchProvider:  researchProvider,
		inferenceProvider: inferenceProvider,
		taskCh:            make(chan inference.InferenceTask),
	}

	go research.inferenceWorker()

	return research
}

func (s *ResearchService) RunFileProcessing(filename, collectionId string, src io.Reader) error {
	err := s.researchProvider.SaveFile(collectionId, filename, src)
	if err != nil {
		s.log.Error("fail save file", slog.String("filename", filename), slog.String("err", err.Error()))
		return err
	}

	id := uuid.New().String()
	filepath := s.cfg.ResearchSavePath + "/" + collectionId + "/" + filename
	err = s.researchProvider.Create(id, collectionId, filepath)
	if err != nil {
		s.log.Error("fail create row in db", slog.String("err", err.Error()))
	}

	go func() {
		// горутина нужна, чтобы не блокировать HTTP-вызов
		// она сохранит контекст и встанет дожидаться очереди на отправку
		s.taskCh <- inference.InferenceTask{
			ResearchId: id,
			Filepath:   filepath,
		}
	}()

	return nil
}

func (s *ResearchService) inferenceWorker() {
	s.log.Info("start inference worker")

	for t := range s.taskCh {
		s.log.Info("start inference", slog.String("filepath", t.Filepath))
		responseCh := make(chan inference.InferenceResponse)

		go func() {
			defer close(responseCh)
			if err := s.inferenceProvider.DoInference(responseCh, t.Filepath); err != nil {
				s.log.Error("inference error", slog.String("err", err.Error()))
			}
		}()

		var inferenceResponse inference.InferenceResponse

		for r := range responseCh {
			s.log.Info("inference",
				slog.Bool("done", r.Done),
				slog.String("step", r.Step),
				slog.Uint64("percent", uint64(r.Percent)),
				slog.Float64("ProbabilityOfPathology", float64(r.ProbabilityOfPathology)))
			inferenceResponse = r
		}

		s.log.Info("finish inference")
		finishedAt := time.Now().UTC()

		if err := s.researchProvider.WriteInferenceResult(t.ResearchId, inferenceResponse.ProbabilityOfPathology, finishedAt); err != nil {
			s.log.Error("fail write inference result in db", slog.String("err", err.Error()))
		} else {
			s.log.Debug("inference result writed to DB")
		}
	}

	s.log.Warn("finish inference worker")
}
