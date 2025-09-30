package research

import (
	"apiservice/internal/domain/research"
	"archive/zip"
	"log/slog"
	"strings"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func (s *ResearchService) readMetadatas(reader *zip.ReadCloser, setTask func(research.Metadata)) error {
	uniqMetadatas := make(map[string]research.Metadata)

	for _, f := range reader.File {
		info := f.FileInfo()
		if info.IsDir() {
			continue // избегаем вложенности для кода ниже
		}

		bytesToRead := info.Size()
		rc, err := f.Open()
		if err != nil {
			s.log.Warn("can't open file in ZIP", slog.String("err", err.Error()), slog.String("filename", f.Name))
			continue
		}

		d, err := dicom.Parse(rc, bytesToRead, nil)
		rc.Close()
		if err != nil {
			s.log.Warn("can't parse DICOM", slog.String("err", err.Error()), slog.String("filename", f.Name))
			continue
		}

		studyIdElem, err := d.FindElementByTag(tag.StudyInstanceUID)
		if err != nil {
			s.log.Warn("not found studyIdElem in DICOM", slog.String("err", err.Error()), slog.String("filename", info.Name()))
		}

		seriesIdElem, err := d.FindElementByTag(tag.SeriesInstanceUID)
		if err != nil {
			s.log.Warn("not found seriesIdElem in DICOM", slog.String("err", err.Error()), slog.String("filename", info.Name()))
		}

		if studyIdElem != nil && seriesIdElem != nil {
			studyId := strings.Trim(studyIdElem.Value.String(), "[]")
			seriesId := strings.Trim(seriesIdElem.Value.String(), "[]")
			key := studyId + "_" + seriesId

			filesCount := 1 // текущий прочитанный файл входит в общее количество файлов
			metadata, ok := uniqMetadatas[key]
			if ok {
				filesCount = metadata.FilesCount + 1
			} else {
				s.log.Info("find uniq metadata; create inference task", slog.String("key", key))

				go setTask(metadata)
			}

			uniqMetadatas[key] = research.Metadata{
				StudyId:    studyId,
				SeriesId:   seriesId,
				FilesCount: filesCount,
			}
		}
	}

	s.log.Info("metadatas count", slog.Int("count", len(uniqMetadatas)))

	if len(uniqMetadatas) == 0 {
		return research.ErrNotFoundMetadata
	}

	metadatas := make([]research.Metadata, 0, 1) // чаще всего в архиве будет 1 серия
	for _, v := range uniqMetadatas {
		metadatas = append(metadatas, v)
	}

	return nil
}
