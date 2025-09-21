package research

import (
	"apiservice/internal/domain/research"
	"database/sql"
)

type ResearchStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *ResearchStorage {
	return &ResearchStorage{
		db: db,
	}
}

func (s *ResearchStorage) Create(id string, pathologyLevel float32) (research.Research, error) {
	q := `INSERT INTO researches (id,pathology_level) VALUES ($1,$2)
				RETURNING id,num,title,pathology_level,created_at`
	var (
		res   research.Research
		title sql.NullString
	)

	row := s.db.QueryRow(q, id, pathologyLevel)

	if err := row.Scan(&res.Id, &res.Num, &title, &res.PathologyLevel, &res.CreatedAt); err != nil {
		return research.Research{}, err
	}

	if title.Valid {
		res.Title = title.String
	}

	return res, nil
}
