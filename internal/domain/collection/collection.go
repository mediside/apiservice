package collection

import (
	"apiservice/internal/domain/research"
	"time"
)

type Collection struct {
	Id             string    `json:"id"`
	Num            uint      `json:"num"`
	Title          string    `json:"title"`
	Folder         string    `json:"folder"`
	PathologyLevel float32   `json:"pathology_level"`
	CreatedAt      time.Time `json:"created_at"`
}

type Update struct {
	Title          *string  `json:"title,omitempty"`
	PathologyLevel *float32 `json:"pathology_level,omitempty"`
}

type CollectionWithResearches struct {
	Id             string                    `json:"id"`
	Num            uint                      `json:"num"`
	Title          string                    `json:"title"`
	Folder         string                    `json:"folder"`
	PathologyLevel float32                   `json:"pathology_level"`
	CreatedAt      time.Time                 `json:"created_at"`
	Researches     []research.ResearchResult `json:"researches"`
}
