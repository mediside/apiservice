package research

type Metadata struct {
	StudyId    string `json:"study_id"`
	SeriesId   string `json:"series_id"`
	FilesCount int    `json:"files_count"`
}

func (m *Metadata) IsZero() bool {
	return m.FilesCount == 0 && m.StudyId == "" && m.SeriesId == ""
}
