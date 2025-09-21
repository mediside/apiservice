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

	row := s.db.QueryRow(q, id, pathologyLevel)

	return s.scanRow(row.Scan)
}

func (s *ResearchStorage) Delete(id string) error {
	q := "DELETE FROM researches WHERE id = $1"
	_, err := s.db.Exec(q, id)
	return err
}

func (s *ResearchStorage) List() ([]research.Research, error) {
	q := "SELECT * FROM researches"
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []research.Research // TODO: оптимизировать capacity через SELECT COUNT(*) FROM table
	// res := make([]research.Research, 0)

	for rows.Next() {
		r, err := s.scanRow(rows.Scan)
		if err != nil {
			return nil, err
		}

		res = append(res, r)
	}

	return res, nil
}

func (s *ResearchStorage) GetOne(id string) (research.Research, error) {
	q := "SELECT * FROM researches WHERE id = $1"
	row := s.db.QueryRow(q, id)

	return s.scanRow(row.Scan)
}

func (s *ResearchStorage) scanRow(scanFn func(dest ...any) error) (research.Research, error) {
	var (
		res   research.Research
		title sql.NullString
	)

	if err := scanFn(&res.Id, &res.Num, &title, &res.PathologyLevel, &res.CreatedAt); err != nil {
		return research.Research{}, err
	}

	if title.Valid {
		res.Title = title.String
	}

	return res, nil
}
