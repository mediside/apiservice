package research

import (
	"apiservice/internal/domain/research"
	"database/sql"
	"io"
	"os"
	"time"
)

type ResearchStorage struct {
	researchSavePath string
	db               *sql.DB
}

func New(researchSavePath string, db *sql.DB) *ResearchStorage {
	return &ResearchStorage{
		researchSavePath: researchSavePath,
		db:               db,
	}
}

func (s *ResearchStorage) SaveFile(subfolder, filename string, src io.Reader) error {
	folderpath := s.researchSavePath + "/" + subfolder
	if err := os.MkdirAll(folderpath, os.ModePerm); err != nil {
		return err
	}

	filepath := folderpath + "/" + filename
	if _, err := os.Stat(filepath); err == nil {
		return research.ErrFileAlreadyExists
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return err
	}

	return nil
}

func (s *ResearchStorage) Create(id, collectionId, filepath string) error {
	q := "INSERT INTO researches (id,collection_id,file_path) VALUES ($1,$2,$3)"
	_, err := s.db.Exec(q, id, collectionId, filepath)
	return err
}

func (s *ResearchStorage) WriteInferenceResult(id string, probabilityOfPathology float32) error {
	q := "UPDATE researches SET probability_of_pathology = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, probabilityOfPathology)
	return err
}

func (s *ResearchStorage) WriteInferenceError(id, inferenceErr string) error {
	q := "UPDATE researches SET inference_error = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, inferenceErr)
	return err
}

func (s *ResearchStorage) WriteInferenceFinishTime(id string, finishedAt time.Time) error {
	q := "UPDATE researches SET processing_finished_at = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, finishedAt)
	return err
}

func (s *ResearchStorage) WriteMetadata(id string, metadata research.ResearchMetadata, size int64) error {
	q := "UPDATE researches SET study_id = $2, series_id = $3, files_count = $4, archive_size = $5 WHERE id = $1"
	_, err := s.db.Exec(q, id, metadata.StudyId, metadata.SeriesId, metadata.FilesCount, size)
	return err
}
