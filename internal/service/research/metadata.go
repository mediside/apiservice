package research

import (
	"apiservice/internal/domain/research"
	"archive/zip"
	"log/slog"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func (s *ResearchService) readMetadata(id string, reader *zip.ReadCloser) (research.ResearchMetadata, error) {
	// нужны: study_uid и series_uid (из DICOM-тегов)
	var (
		studyIdElem  *dicom.Element
		seriesIdElem *dicom.Element
	)

	// последовательно идем по всем файлам. Когда прочитали оба поля - не важно, из каких файлов
	// то дальше файлы не читаем
	for _, f := range reader.File {
		info := f.FileInfo()
		if !info.IsDir() {
			bytesToRead := info.Size()
			rc, err := f.Open()
			if err != nil {
				s.log.Error("can't open file in ZIP", slog.String("err", err.Error()))
				return research.ResearchMetadata{}, err
			}

			d, err := dicom.Parse(rc, bytesToRead, nil)
			rc.Close()
			if err != nil {
				s.log.Error("can't parse DICOM", slog.String("err", err.Error()))
				return research.ResearchMetadata{}, err
			}

			if studyIdElem == nil {
				studyIdElem, err = d.FindElementByTag(tag.StudyInstanceUID)
				if err != nil {
					s.log.Warn("not found studyIdElem", slog.String("err", err.Error()), slog.String("filename", info.Name()))
				}
			}

			if seriesIdElem == nil {
				seriesIdElem, err = d.FindElementByTag(tag.SeriesInstanceUID)
				if err != nil {
					s.log.Warn("not found seriesIdElem", slog.String("err", err.Error()), slog.String("filename", info.Name()))
				}
			}

			if studyIdElem != nil && seriesIdElem != nil {
				break
			}
		}
	}

	if studyIdElem == nil || seriesIdElem == nil {
		s.log.Error("not found study_uid or series_uid", slog.String("id", id))
		return research.ResearchMetadata{}, research.ErrNotFoundMetadata
	}
	studyId := studyIdElem.Value.String()
	seriesId := seriesIdElem.Value.String()
	return research.ResearchMetadata{StudyId: studyId, SeriesId: seriesId, FilesCount: len(reader.File)}, nil
}
