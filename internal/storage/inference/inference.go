package inference

import infGRPC "apiservice/internal/gen/go/inference/inference.v1"

type InferenceStorage struct {
	client infGRPC.InferenceClient
}

func New(client infGRPC.InferenceClient) *InferenceStorage {
	return &InferenceStorage{
		client: client,
	}
}

func DoInference(filepath string) {

}
