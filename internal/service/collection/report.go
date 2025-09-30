package collection

import (
	"bytes"
	"log/slog"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func (s *CollectionService) CreateReport(id string) (*bytes.Buffer, error) {
	col, err := s.collectionProvider.GetOne(id)
	if err != nil {
		s.log.Error("get one collection", slog.String("err", err.Error()))
		return nil, err
	}

	rs, err := s.researchProvider.List(id)
	if err != nil {
		s.log.Error("list researches", slog.String("err", err.Error()))
		return nil, err
	}

	pathologyLevel := col.PathologyLevel

	f := excelize.NewFile()
	sheet := "Sheet1"

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		s.log.Error("unable create style for XLSX table header")
		return nil, err
	}

	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}

	f.SetRowStyle(sheet, 1, 1, headerStyle)

	// заголовок
	f.SetCellValue(sheet, "A1", "path_to_study")
	f.SetCellValue(sheet, "B1", "study_uid")
	f.SetCellValue(sheet, "C1", "series_uid")
	f.SetCellValue(sheet, "D1", "probability_of_pathology")
	f.SetCellValue(sheet, "E1", "pathology")
	f.SetCellValue(sheet, "F1", "processing_status")
	f.SetCellValue(sheet, "G1", "time_of_processing")
	// дополнительные поля (по собственной инициативе)
	f.SetCellValue(sheet, "H1", "processing_error") // причину ошибки знать полезно

	rowIndex := 2
	for _, r := range rs {
		if r.ArchiveCorrupt {
			continue // поврежденный архив в отчет не включаем
		}
		iStr := strconv.Itoa(rowIndex)
		rowIndex++

		f.SetCellValue(sheet, "A"+iStr, r.Filepath)
		f.SetCellValue(sheet, "B"+iStr, r.Metadata.StudyId)
		f.SetCellValue(sheet, "C"+iStr, r.Metadata.SeriesId)

		f.SetCellValue(sheet, "D"+iStr, r.ProbabilityOfPathology)
		isPathology := uint8(0)
		if r.ProbabilityOfPathology > pathologyLevel {
			isPathology = 1
		}
		f.SetCellValue(sheet, "E"+iStr, isPathology)

		processingStatus := "Success"
		if r.InferenceError != "" {
			processingStatus = "Failure"
		}
		f.SetCellValue(sheet, "F"+iStr, processingStatus)

		diff := r.ProcessingFinishedAt.Sub(r.ProcessingStartedAt)
		f.SetCellValue(sheet, "G"+iStr, int(diff.Seconds()))

		f.SetCellValue(sheet, "H"+iStr, r.InferenceError)
	}

	f.SetActiveSheet(index)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
