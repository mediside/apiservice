package collection

import "time"

type Collection struct {
	Id             string    `json:"id"`
	Num            uint      `json:"num"`
	Title          string    `json:"title"`
	PathologyLevel float32   `json:"pathology_level"`
	CreatedAt      time.Time `json:"created_at"`
}
