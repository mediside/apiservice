package research

import (
	"apiservice/internal/domain/inference"
	"apiservice/internal/domain/research"
	"archive/zip"
	"log/slog"
	"os"
	filepathLib "path/filepath"

	"github.com/google/uuid"
)

func (s *Service) processing(filename, collectionId string) {
	filepath := s.cfg.ResearchSavePath + "/" + collectionId + "/" + filename

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		s.log.Error("can't get file stat", slog.String("filepath", filepath), slog.String("err", err.Error()))
		return
	}
	size := fileInfo.Size()

	// перед отправкой в инференс кратко смотрим, не битый ли архив
	reader, err := zip.OpenReader(filepath)
	if err != nil {
		s.log.Error("can't open ZIP", slog.String("err", err.Error()))
		id := uuid.New().String()
		if err := s.researchProvider.Create(id, collectionId, filepath, size, true, research.Metadata{}); err != nil {
			s.log.Error("fail create corrupted research in db", slog.String("err", err.Error()), slog.String("id", id))
		} else {
			s.updateCh <- research.ResearchUpdate{
				Id:             id,
				CollectionId:   collectionId,
				Filepath:       filepath,
				Filename:       filepathLib.Base(filepath),
				Size:           size,
				ArchiveCorrupt: true,
			}
		}

		if err := s.researchProvider.DeleteSingleFile(filepath); err != nil {
			s.log.Warn("can't delete file", slog.String("err", err.Error()), slog.String("filepath", filepath))
		} // оставляем запись в БД, но удаляем поврежденный файл

		return // если не смогли сами прочитать архив, то не даем задачу на инференс
	}
	defer reader.Close()

	if err := s.readMetadatas(reader, func(metadata research.Metadata) {
		researchId := uuid.New().String()
		if err := s.researchProvider.Create(researchId, collectionId, filepath, size, false, metadata); err != nil {
			s.log.Warn("fail create research with metedata in db", slog.String("err", err.Error()), slog.String("id", researchId))
			return
		}

		s.updateCh <- research.ResearchUpdate{
			Id:           researchId,
			CollectionId: collectionId,
			Filepath:     filepath,
			Filename:     filepathLib.Base(filepath),
			Size:         size,
			Metadata:     metadata,
		}

		s.taskCh <- inference.InferenceTask{
			ResearchId:   researchId,
			CollectionId: collectionId,
			Filepath:     filepath,
			Size:         size,
			Metadata:     metadata,
		}
	}); err != nil {
		return
	}
}
