package tmf

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
)

func GetModelFrom3mfFile(file3mf string) string {
	zipBytes, err := os.ReadFile(file3mf)
	if err != nil {
		return ""
	}
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return ""
	}

	for _, file := range zipReader.File {
		if file.Name == "Metadata/Slic3r_PE.config" {
			// Öffne die gewünschte Datei
			fileReader, err := file.Open()
			if err != nil {
				return ""
			}
			defer fileReader.Close()

			// Lese die Datei in ein Byte-Array
			var buffer bytes.Buffer
			if _, err = io.Copy(&buffer, fileReader); err != nil {
				return ""
			}

			return string(buffer.Bytes())
		}
	}

	return ""
}
