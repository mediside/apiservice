package research

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/inference"
	"apiservice/internal/domain/research"
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

type InferenceProvider interface {
	DoInference(responseCh chan<- inference.InferenceResponse, filepath string) error
}

type ResearchProvider interface {
	SaveFile(collectionId, filename string, src io.Reader) error
	Create(id, collectionId, filepath string) error
	Delete(id string) error
	WriteInferenceResult(id string, probabilityOfPathology float32) error
	WriteInferenceError(id, inferenceErr string) error
	WriteInferenceFinishTime(id string, finishedAt time.Time) error
	WriteInferenceStartTime(id string, startedAt time.Time) error
	WriteMetadata(id string, metadata research.ResearchMetadata, size int64) error
	MarkCorrupted(id string) error
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

func (s *ResearchService) processing(filename, collectionId string) {
	id := uuid.New().String()
	filepath := s.cfg.ResearchSavePath + "/" + collectionId + "/" + filename
	err := s.researchProvider.Create(id, collectionId, filepath)
	if err != nil {
		s.log.Error("fail create row in db", slog.String("err", err.Error()))
	}

	// перед отправкой в инференс кратко смотрим, не битый ли архив
	reader, err := zip.OpenReader(filepath)
	if err != nil {
		s.log.Error("can't open ZIP", slog.String("err", err.Error()))
		s.researchProvider.MarkCorrupted(id)
		return // если не смогли сами прочитать архив, то не даем задачу на инференс
	}
	defer reader.Close()

	go func() {
		// ожидание инференса в очереди не блокирует запись информации об исследовании
		s.taskCh <- inference.InferenceTask{
			ResearchId: id,
			Filepath:   filepath,
		}
	}()

	metadata, err := s.readMetadata(id, reader)
	if err != nil {
		return
	}

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		s.log.Error("can't get file stat", slog.String("filepath", filepath), slog.String("err", err.Error()))
		return
	}
	size := fileInfo.Size()

	err = s.researchProvider.WriteMetadata(id, metadata, size)
	if err != nil {
		s.log.Error("can't write metadata", slog.String("id", id), slog.String("err", err.Error()))
	}
}

func (s *ResearchService) readMetadata(id string, reader *zip.ReadCloser) (research.ResearchMetadata, error) {
	// нужны: study_uid и series_uid (из DICOM-тегов)
	var (
		studyIdElem  *dicom.Element
		seriesIdElem *dicom.Element
	)

	// последовательно идем по всем файлам. Когда прочитали оба поля - не важно, из каких файлов
	// то дальше файлы не читаем
	for _, f := range reader.File {
		info := f.FileInfo()
		if !info.IsDir() {
			bytesToRead := info.Size()
			rc, err := f.Open()
			if err != nil {
				s.log.Error("can't open file in ZIP", slog.String("err", err.Error()))
				return research.ResearchMetadata{}, err
			}

			d, err := dicom.Parse(rc, bytesToRead, nil)
			rc.Close()
			if err != nil {
				s.log.Error("can't parse DICOM", slog.String("err", err.Error()))
				return research.ResearchMetadata{}, err
			}

			if studyIdElem == nil {
				studyIdElem, err = d.FindElementByTag(tag.StudyInstanceUID)
				if err != nil {
					s.log.Warn("not found studyIdElem", slog.String("err", err.Error()), slog.String("filename", info.Name()))
				}
			}

			if seriesIdElem == nil {
				seriesIdElem, err = d.FindElementByTag(tag.SeriesInstanceUID)
				if err != nil {
					s.log.Warn("not found seriesIdElem", slog.String("err", err.Error()), slog.String("filename", info.Name()))
				}
			}

			if studyIdElem != nil && seriesIdElem != nil {
				break
			}
		}
	}

	if studyIdElem == nil || seriesIdElem == nil {
		s.log.Error("not found study_uid or series_uid", slog.String("id", id))
		return research.ResearchMetadata{}, research.ErrNotFoundMetadata
	}
	studyId := studyIdElem.Value.String()
	seriesId := seriesIdElem.Value.String()
	return research.ResearchMetadata{StudyId: studyId, SeriesId: seriesId, FilesCount: len(reader.File)}, nil
}

func (s *ResearchService) inferenceWorker() {
	s.log.Info("start inference worker")

	for t := range s.taskCh {
		s.log.Info("start inference", slog.String("filepath", t.Filepath))
		startedAt := time.Now().UTC()
		if err := s.researchProvider.WriteInferenceStartTime(t.ResearchId, startedAt); err != nil {
			s.log.Error("fail write inference start time in db", slog.String("err", err.Error()))
		}

		responseCh := make(chan inference.InferenceResponse)

		go func() {
			defer close(responseCh)
			if err := s.inferenceProvider.DoInference(responseCh, t.Filepath); err != nil {
				s.log.Warn("inference error", slog.String("err", err.Error()))
				if e := s.researchProvider.WriteInferenceError(t.ResearchId, err.Error()); e != nil {
					s.log.Error("fail write inference error in db", slog.String("err", err.Error()))
				} else {
					s.log.Debug("inference error writed to DB")
				}

				s.inferenceCh <- inference.InferenceProgress{
					Done:       true, // инференс закончен, но с ошибкой
					ResearchId: t.ResearchId,
					SeriesId:   "", // TODO
					StudyId:    "", // TODO
					Err:        err.Error(),
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
				SeriesId:               "", // TODO
				StudyId:                "", // TODO
			}
		}

		finishedAt := time.Now().UTC()
		if err := s.researchProvider.WriteInferenceFinishTime(t.ResearchId, finishedAt); err != nil {
			s.log.Error("fail write inference finish time in db", slog.String("err", err.Error()))
		} else {
			s.log.Debug("inference finish time writed to DB")
		}
		s.log.Info("finish inference")

		if inferenceResponse.Done {
			s.log.Info("inference success", slog.String("ResearchId", t.ResearchId))
			if err := s.researchProvider.WriteInferenceResult(t.ResearchId, inferenceResponse.ProbabilityOfPathology); err != nil {
				s.log.Error("fail write inference result in db", slog.String("err", err.Error()))
			} else {
				s.log.Debug("inference result writed to DB")
			}
		}
	}

	s.log.Warn("finish inference worker")
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
