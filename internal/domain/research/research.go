package research

import (
	"errors"
	"time"
)

var ErrFileAlreadyExists = errors.New("file already exists in folder")
var ErrNotFoundMetadata = errors.New("metadata not found in DICOM research")

type ResearchMetadata struct {
	StudyId    string
	SeriesId   string
	FilesCount int
}

type Research struct {
	Id                     string
	Filepath               string
	Size                   int64
	Assessment             string
	ArchiveCorrupt         bool
	ProbabilityOfPathology float32
	CreatedAt              time.Time
	ProcessingStartedAt    time.Time
	ProcessingFinishedAt   time.Time
	InferenceError         string
	Metadata               Metadata
}

type ResearchResult struct {
	Id                     string    `json:"id"`
	Filepath               string    `json:"filepath"`
	Filename               string    `json:"filename"`
	Size                   int64     `json:"size"`
	Assessment             string    `json:"assessment,omitempty"`
	ArchiveCorrupt         bool      `json:"archive_corrupt"`
	ProbabilityOfPathology float32   `json:"probability_of_pathology,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	ProcessingStartedAt    time.Time `json:"processing_started_at,omitzero"`
	ProcessingFinishedAt   time.Time `json:"processing_finished_at,omitzero"`
	ProcessingDuration     int64     `json:"processing_duration,omitzero"` // в миллисекундах
	InferenceError         string    `json:"inference_error,omitempty"`
	Metadata               Metadata  `json:"metadata,omitzero"`
}
