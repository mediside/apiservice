package inference

import (
	"apiservice/internal/domain/inference"
	infGRPC "apiservice/internal/gen/go/inference/inference.v1"
	"context"
	"time"
)

const eps = 0.0001

type Storage struct {
	client infGRPC.InferenceClient
}

func New(client infGRPC.InferenceClient) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) DoInference(responseCh chan<- inference.InferenceResponse, filepath, studyId, seriesId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	stream, err := s.client.DoInference(ctx, &infGRPC.InferenceRequest{FilePath: filepath, StudyId: studyId, SeriesId: seriesId})
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
			prob := payload.Result.ProbabilityOfPathology
			// Это поле в структуре ResearchResult помечено как omitempty. Нельзя допустить, чтобы значение было нулевым.
			// В таком случае оно не будет отправлено в интерфейс вообще. Чтобы этого избежать, нужно присвоить небольшое значение
			if prob < eps {
				prob = eps
			}
			responseCh <- inference.InferenceResponse{
				ProbabilityOfPathology: prob,
				Done:                   true,
			}
			return nil // прерываем цикл
		default:
			return inference.ErrGrpcUnknown
		}
	}
}
