package research

import "errors"

var ErrFileAlreadyExists = errors.New("file already exists in folder")
var ErrNotFoundMetadata = errors.New("metadata not found in DICOM research")

type ResearchMetadata struct {
	StudyId    string
	SeriesId   string
	FilesCount int
}
