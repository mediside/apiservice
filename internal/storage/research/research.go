package research

import (
	"database/sql"
	"io"
	"os"
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

func (s *ResearchStorage) SaveFile(filename string, src io.Reader) error {
	out, err := os.Create("./researches/" + filename)
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
	q := `INSERT INTO collections (id,collection_id,file_path) VALUES ($1,$2,$3)`
	_, err := s.db.Exec(q, id, collectionId, filepath)
	return err
}
