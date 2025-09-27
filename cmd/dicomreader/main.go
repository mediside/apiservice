package main

import (
	"archive/zip"
	"fmt"
	"log"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func main() {
	filepath := "/home/m/dev/hackathon/apiservice/researches/45e15946-9389-4919-8c4c-a8cc7fd46a08/pneumonia_anon.zip"
	r, err := zip.OpenReader(filepath)
	if err != nil {
		log.Fatalf("can't open ZIP: %s", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			bytesToRead := f.FileInfo().Size()
			rc, err := f.Open()
			if err != nil {
				log.Println(err)
				continue
			}

			// Парсим DICOM из потока в памяти
			d, err := dicom.Parse(rc, bytesToRead, nil)
			rc.Close()
			if err != nil {
				log.Println("Error parsing DICOM:", err)
				continue
			}

			fmt.Println("Parsed DICOM file:", f.Name)

			patientNameElem, err := d.FindElementByTag(tag.SeriesInstanceUID)
			if err != nil {
				fmt.Println("tag not found:", err)
			} else {
				fmt.Println("tag:", patientNameElem.Value.String())
				fmt.Println("str:", patientNameElem.ValueRepresentation)
			}
		}

	}
}
