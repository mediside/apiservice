package collection

import (
	"apiservice/internal/domain/research"
	"time"
)

type Collection struct {
	Id             string    `json:"id"`
	Num            uint      `json:"num"`
	Title          string    `json:"title"`
	PathologyLevel float32   `json:"pathology_level"`
	CreatedAt      time.Time `json:"created_at"`
}

type CollectionWithResearches struct {
	Id             string                    `json:"id"`
	Num            uint                      `json:"num"`
	Title          string                    `json:"title"`
	PathologyLevel float32                   `json:"pathology_level"`
	CreatedAt      time.Time                 `json:"created_at"`
	Researches     []research.ResearchResult `json:"researches"`
}
