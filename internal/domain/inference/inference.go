package inference

import "errors"

var ErrGrpcUnknown = errors.New("unknown message in gRPC")

type InferenceTask struct {
	ResearchId string
	Filepath   string
}

type InferenceResponse struct {
	Percent                uint
	Step                   string
	ProbabilityOfPathology float32
	Done                   bool
}
