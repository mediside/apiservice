package inference

import (
	"apiservice/internal/domain/inference"
	infGRPC "apiservice/internal/gen/go/inference/inference.v1"
	"context"
	"time"
)

type InferenceStorage struct {
	client infGRPC.InferenceClient
}

func New(client infGRPC.InferenceClient) *InferenceStorage {
	return &InferenceStorage{
		client: client,
	}
}

func (s *InferenceStorage) DoInference(responseCh chan<- inference.InferenceResponse, filepath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	stream, err := s.client.DoInference(ctx, &infGRPC.InferenceRequest{FilePath: filepath})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		switch payload := resp.Payload.(type) {
		case *infGRPC.InferenceResponse_Progress:
			responseCh <- inference.InferenceResponse{
				Percent: uint(payload.Progress.Percent),
				Step:    payload.Progress.Step,
			}
		case *infGRPC.InferenceResponse_Result:
			responseCh <- inference.InferenceResponse{
				ProbabilityOfPathology: payload.Result.ProbabilityOfPathology,
				Done:                   true,
			}
			return nil // прерываем цикл
		default:
			return inference.ErrGrpcUnknown
		}
	}
}
