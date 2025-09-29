package inference

import "errors"

var ErrGrpcUnknown = errors.New("unknown message in gRPC")

type InferenceTask struct {
	ResearchId   string
	CollectionId string
	Filepath     string
	StudyId      string // в целом, достаточно только одного Id, но можно оставить оба
	SeriesId     string // для поддержки нескольких разных исследований в одном архиве
}

type InferenceResponse struct {
	Percent                uint
	Step                   string
	ProbabilityOfPathology float32
	Done                   bool
}

type InferenceProgress struct {
	Percent                uint    `json:"percent"`
	Step                   string  `json:"step"`
	ProbabilityOfPathology float32 `json:"probability_of_pathology"`
	Done                   bool    `json:"done"`
	ResearchId             string  `json:"research_id"`
	StudyId                string  `json:"study_id"`
	SeriesId               string  `json:"series_id"`
	Err                    string  `json:"err"`
}
