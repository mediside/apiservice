package research

import (
	"io"
	"os"
)

type ResearchStorage struct {
	researchSavePath string
}

func New(researchSavePath string) *ResearchStorage {
	return &ResearchStorage{
		researchSavePath: researchSavePath,
	}
}

func (s *ResearchStorage) SaveFile(filename string, src io.Reader) error {
	out, err := os.Create("./researches/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return err
	}

	return nil
}
