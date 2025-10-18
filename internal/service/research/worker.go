package research

import (
	"apiservice/internal/domain/inference"
	"apiservice/internal/domain/research"
	"log/slog"
	filepathLib "path/filepath"
	"time"
)

func (s *Service) inferenceWorker() {
	s.log.Info("start inference worker")

	for t := range s.taskCh {

		s.updateCh <- research.ResearchUpdate{
			Id:           t.ResearchId,
			CollectionId: t.CollectionId,
			Filepath:     t.Filepath,
			Filename:     filepathLib.Base(t.Filepath),
			Size:         t.Size,
			Metadata:     t.Metadata,
		}

		s.log.Info("start inference", slog.String("filepath", t.Filepath))
		startedAt := time.Now().UTC()
		if err := s.researchProvider.WriteInferenceStartTime(t.ResearchId, startedAt); err != nil {
			s.log.Error("fail write inference start time in db", slog.String("err", err.Error()))
		} else {
			s.updateCh <- research.ResearchUpdate{
				Id:                  t.ResearchId,
				CollectionId:        t.CollectionId,
				ProcessingStartedAt: startedAt,
			}
		}

		responseCh := make(chan inference.InferenceResponse)

		go func() {
			defer close(responseCh)
			if err := s.inferenceProvider.DoInference(responseCh, t.Filepath, t.Metadata.StudyId, t.Metadata.SeriesId); err != nil {
				s.log.Warn("inference error", slog.String("err", err.Error()))
				if e := s.researchProvider.WriteInferenceError(t.ResearchId, err.Error()); e != nil {
					s.log.Error("fail write inference error in db", slog.String("err", err.Error()))
				} else {
					s.updateCh <- research.ResearchUpdate{
						Id:             t.ResearchId,
						CollectionId:   t.CollectionId,
						InferenceError: err.Error(),
					}
					s.log.Debug("inference error writed to DB")
				}

				s.inferenceCh <- inference.InferenceProgress{
					Done:         true, // инференс закончен, но с ошибкой
					ResearchId:   t.ResearchId,
					CollectionId: t.CollectionId,
					SeriesId:     t.Metadata.SeriesId,
					StudyId:      t.Metadata.StudyId,
					Err:          err.Error(),
				}
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

			s.inferenceCh <- inference.InferenceProgress{
				Percent:                r.Percent,
				Step:                   r.Step,
				ProbabilityOfPathology: r.ProbabilityOfPathology,
				Done:                   r.Done,
				ResearchId:             t.ResearchId,
				CollectionId:           t.CollectionId,
				SeriesId:               t.Metadata.SeriesId,
				StudyId:                t.Metadata.StudyId,
			}
		}

		if _, ok := s.counts[t.Filepath]; ok {
			// В теории возможен сценарий, когда в одном файле 2 исследования, на первом инференс прошел очень быстро, а второе исследование сервис еще не успел найти
			// В таком случае счетчик уменьшится до нуля, файл будет удален до инференса второго исследования
			// Но на практике такая ситуация маловероятна
			s.mu.Lock()
			s.counts[t.Filepath] -= 1
			s.mu.Unlock()
			if s.counts[t.Filepath] == 0 {
				if err := s.researchProvider.DeleteSingleFile(t.Filepath); err != nil {
					s.log.Warn("can't delete file", slog.String("err", err.Error()), slog.String("filepath", t.Filepath))
				}
			}
		} else {
			s.log.Warn("not found filepath in counts", slog.String("filepath", t.Filepath))
		}

		finishedAt := time.Now().UTC()
		if err := s.researchProvider.WriteInferenceFinishTime(t.ResearchId, finishedAt); err != nil {
			s.log.Error("fail write inference finish time in db", slog.String("err", err.Error()))
		} else {
			s.updateCh <- research.ResearchUpdate{
				Id:                   t.ResearchId,
				CollectionId:         t.CollectionId,
				ProcessingFinishedAt: finishedAt,
				ProcessingDuration:   int64(finishedAt.Sub(startedAt).Seconds()),
			}
			s.log.Debug("inference finish time writed to DB")
		}
		s.log.Info("finish inference")

		if inferenceResponse.Done {
			s.log.Info("inference success", slog.String("ResearchId", t.ResearchId))
			if err := s.researchProvider.WriteInferenceResult(t.ResearchId, inferenceResponse.ProbabilityOfPathology); err != nil {
				s.log.Error("fail write inference result in db", slog.String("err", err.Error()))
			} else {
				s.updateCh <- research.ResearchUpdate{
					Id:                     t.ResearchId,
					CollectionId:           t.CollectionId,
					ProbabilityOfPathology: inferenceResponse.ProbabilityOfPathology,
				}
				s.log.Debug("inference result writed to DB")
			}
		}
	}

	s.log.Warn("finish inference worker")
}
