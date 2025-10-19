package collection

import (
	"apiservice/internal/domain/collection"
	"database/sql"
	"os"
)

type Storage struct {
	savePath string
	db       *sql.DB
}

func New(savePath string, db *sql.DB) *Storage {
	return &Storage{
		savePath: savePath,
		db:       db,
	}
}

func (s *Storage) Create(id string, pathologyLevel float32) (collection.Collection, error) {
	// TODO: разбить на подпакеты работу с файловой системой и базой данных
	folderpath := s.savePath + "/" + id
	if _, err := os.Stat(folderpath); os.IsNotExist(err) {
		err := os.MkdirAll(folderpath, 0755)
		if err != nil {
			return collection.Collection{}, err
		}
	}

	q := `INSERT INTO collections (id,pathology_level) VALUES ($1,$2)
				RETURNING id,num,title,pathology_level,created_at`

	row := s.db.QueryRow(q, id, pathologyLevel)

	return s.scanRow(row.Scan)
}

func (s *Storage) Delete(id string) error {
	q := "DELETE FROM collections WHERE id = $1"
	_, err := s.db.Exec(q, id)
	return err
}

func (s *Storage) List() ([]collection.Collection, error) {
	q := "SELECT * FROM collections"
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []collection.Collection // TODO: оптимизировать capacity через SELECT COUNT(*) FROM table

	for rows.Next() {
		r, err := s.scanRow(rows.Scan)
		if err != nil {
			return nil, err
		}

		res = append(res, r)
	}

	return res, nil
}

func (s *Storage) CheckExists(id string) (bool, error) {
	q := "SELECT EXISTS(SELECT 1 FROM collections WHERE id = $1)"
	row := s.db.QueryRow(q, id)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *Storage) GetOne(id string) (collection.Collection, error) {
	q := "SELECT * FROM collections WHERE id = $1"
	row := s.db.QueryRow(q, id)
	return s.scanRow(row.Scan)
}

func (s *Storage) scanRow(scanFn func(dest ...any) error) (collection.Collection, error) {
	var (
		res   collection.Collection
		title sql.NullString
	)

	if err := scanFn(&res.Id, &res.Num, &title, &res.PathologyLevel, &res.CreatedAt); err != nil {
		return collection.Collection{}, err
	}

	if title.Valid {
		res.Title = title.String
	}

	return res, nil
}

func (s *Storage) WritePathologyLevel(id string, pathologyLevel float32) error {
	q := "UPDATE collections SET pathology_level = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, pathologyLevel)
	return err
}

func (s *Storage) WriteTitle(id string, title string) error {
	q := "UPDATE collections SET title = $2 WHERE id = $1"
	_, err := s.db.Exec(q, id, title)
	return err
}
