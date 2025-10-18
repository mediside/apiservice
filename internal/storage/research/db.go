package research

import (
	"apiservice/internal/domain/research"
	"database/sql"
	"time"
)

func (s *Storage) Create(id, collectionId, filepath string, size int64, archiveCorrupt bool, metadata research.Metadata) error {
	var err error = nil
	if archiveCorrupt {
		q := "INSERT INTO researches (id,collection_id,file_path,archive_size,archive_corrupt) VALUES ($1,$2,$3,$4,$5)"
		_, err = s.db.Exec(q, id, collectionId, filepath, size, true)
	} else {
		q := "INSERT INTO researches (id,collection_id,file_path,archive_size,study_id,series_id,files_count) VALUES ($1,$2,$3,$4,$5,$6,$7)"
		_, err = s.db.Exec(q, id, collectionId, filepath, size, metadata.StudyId, metadata.SeriesId, metadata.FilesCount)
	}

	return err
}

func (s *Storage) WriteInferenceResult(id string, probabilityOfPathology float32) error {
	q := "UPDATE researches SET probability_of_pathology = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, probabilityOfPathology)
	return err
}

func (s *Storage) WriteInferenceError(id, inferenceErr string) error {
	q := "UPDATE researches SET inference_error = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, inferenceErr)
	return err
}

func (s *Storage) WriteInferenceFinishTime(id string, finishedAt time.Time) error {
	q := "UPDATE researches SET processing_finished_at = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, finishedAt)
	return err
}

func (s *Storage) WriteInferenceStartTime(id string, startedAt time.Time) error {
	q := "UPDATE researches SET processing_started_at = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, startedAt)
	return err
}

func (s *Storage) List(collectionId string) ([]research.Research, error) {
	q := `SELECT r.id, r.file_path, r.archive_size, r.assessment, r.archive_corrupt, r.probability_of_pathology,
				r.created_at, r.processing_started_at, r.processing_finished_at, r.study_id, r.series_id, r.files_count, r.inference_error
				FROM researches as r WHERE collection_id = $1 ORDER BY created_at ASC`
	rows, err := s.db.Query(q, collectionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rs []research.Research

	for rows.Next() {
		var (
			assessment  sql.NullString
			probability sql.NullFloat64
			startedAt   sql.NullTime
			finishedAt  sql.NullTime
			studyId     sql.NullString
			seriesId    sql.NullString
			filesCount  sql.NullInt32
			infErr      sql.NullString
		)
		r := research.Research{}
		err := rows.Scan(&r.Id, &r.Filepath, &r.Size, &assessment, &r.ArchiveCorrupt, &probability,
			&r.CreatedAt, &startedAt, &finishedAt, &studyId, &seriesId, &filesCount, &infErr)
		if err != nil {
			return nil, err
		}

		if studyId.Valid && seriesId.Valid && filesCount.Valid {
			// эти поля устанавливаются совместно
			r.Metadata = research.Metadata{
				SeriesId:   seriesId.String,
				StudyId:    studyId.String,
				FilesCount: int(filesCount.Int32),
			}
		}
		if assessment.Valid {
			r.Assessment = assessment.String
		}
		if probability.Valid {
			r.ProbabilityOfPathology = float32(probability.Float64)
		}
		if startedAt.Valid {
			r.ProcessingStartedAt = startedAt.Time
		}
		if finishedAt.Valid {
			r.ProcessingFinishedAt = finishedAt.Time
		}
		if infErr.Valid {
			r.InferenceError = infErr.String
		}

		rs = append(rs, r)
	}

	return rs, nil
}

func (s *Storage) GetFilepath(id string) (string, error) {
	q := "SELECT file_path FROM researches WHERE id = $1"
	var filepath string
	err := s.db.QueryRow(q, id).Scan(&filepath)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func (s *Storage) DeleteEntry(filepath string) error {
	q := "DELETE FROM researches WHERE file_path = $1"
	_, err := s.db.Exec(q, filepath)
	return err
}

func (s *Storage) CheckExists(collectionId, filename string) (bool, error) {
	filepath := s.researchSavePath + "/" + collectionId + "/" + filename

	q := "SELECT EXISTS(SELECT 1 FROM researches WHERE collection_id = $1 AND file_path = $2)"
	row := s.db.QueryRow(q, collectionId, filepath)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
