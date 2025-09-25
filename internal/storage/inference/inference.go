package inference

import (
	infDom "apiservice/internal/domain/inference"
	infGRPC "apiservice/internal/gen/go/inference/inference.v1"
	"context"
	"fmt"
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

func (s *InferenceStorage) DoInference(filepath string) error {
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
			fmt.Printf("Прогресс: %v\n", payload.Progress) // TODO: channels
		case *infGRPC.InferenceResponse_Result:
			fmt.Printf("Результат: %v\n", payload.Result) // TODO: channels
			return nil                                    // прерываем цикл
		default:
			return infDom.ErrGrpcUnknown // TODO: channels
		}
	}
}
