package research

import (
	"apiservice/internal/domain/inference"
	"archive/zip"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

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
