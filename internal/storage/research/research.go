package research

import (
	"database/sql"
)

type Storage struct {
	researchSavePath string
	db               *sql.DB
}

func New(researchSavePath string, db *sql.DB) *Storage {
	return &Storage{
		researchSavePath: researchSavePath,
		db:               db,
	}
}
