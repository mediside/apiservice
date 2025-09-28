package research

type Metadata struct {
	StudyId    string `json:"study_id"`
	SeriesId   string `json:"series_id"`
	FilesCount uint   `json:"files_count"`
	Size       int64  `json:"size"`
}

func (m *Metadata) IsZero() bool {
	return m.FilesCount == 0 && m.StudyId == "" && m.SeriesId == "" && m.Size == 0
}
