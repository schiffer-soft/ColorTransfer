package tmf

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
)

func GetModelFrom3mf(bytes3mf []byte) (*Model, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(bytes3mf), int64(len(bytes3mf)))
	if err != nil {
		return nil, err
	}

	for _, file := range zipReader.File {
		if file.Name == "3D/3dmodel.model" {
			// Öffne die gewünschte Datei
			fileReader, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fileReader.Close()

			// Lese die Datei in ein Byte-Array
			var buffer bytes.Buffer
			if _, err = io.Copy(&buffer, fileReader); err != nil {
				return nil, err
			}

			bytes := buffer.Bytes()
			var model Model
			if err = xml.Unmarshal(bytes, &model); err != nil {
				return nil, err
			}

			return &model, nil
		}
	}

	return nil, errors.New("Model wurde nicht gefunden")
}

func GetThumbnailFrom3mf(bytes3mf []byte) []byte {
	zipReader, err := zip.NewReader(bytes.NewReader(bytes3mf), int64(len(bytes3mf)))
	if err != nil {
		return nil
	}

	for _, file := range zipReader.File {
		if file.Name == "Metadata/thumbnail.png" {
			// Öffne die gewünschte Datei
			fileReader, err := file.Open()
			if err != nil {
				return nil
			}
			defer fileReader.Close()

			// Lese die Datei in ein Byte-Array
			var buffer bytes.Buffer
			if _, err = io.Copy(&buffer, fileReader); err != nil {
				return nil
			}

			return buffer.Bytes()
		}
	}

	return nil
}
