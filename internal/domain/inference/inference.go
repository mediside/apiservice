package inference

import (
	"apiservice/internal/domain/research"
	"errors"
)

var ErrGrpcUnknown = errors.New("unknown message in gRPC")

type InferenceTask struct {
	CollectionId string
	Filepath     string
	Size         int64
	Metadata     research.Metadata
}

type InferenceResponse struct {
	Percent                uint
	Step                   string
	ProbabilityOfPathology float32
	Done                   bool
}

type InferenceProgress struct {
	ResearchId             string  `json:"research_id"`
	CollectionId           string  `json:"collection_id"`
	Percent                uint    `json:"percent"`
	Step                   string  `json:"step"`
	ProbabilityOfPathology float32 `json:"probability_of_pathology"`
	Done                   bool    `json:"done"`
	StudyId                string  `json:"study_id"`
	SeriesId               string  `json:"series_id"`
	Err                    string  `json:"err"`
}
