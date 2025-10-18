package research

import (
	"apiservice/internal/domain/research"
	"io"
	"os"
)

func (s *Storage) SaveFile(subfolder, filename string, src io.Reader) error {
	folderpath := s.researchSavePath + "/" + subfolder
	if err := os.MkdirAll(folderpath, os.ModePerm); err != nil {
		return err
	}

	filepath := folderpath + "/" + filename
	if _, err := os.Stat(filepath); err == nil {
		return research.ErrFileAlreadyExists
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteSingleFile(filepath string) error {
	return os.Remove(filepath)
}

func (s *Storage) DeleteFiles(subfolder string) error {
	folderpath := s.researchSavePath + "/" + subfolder
	return os.RemoveAll(folderpath)
}
