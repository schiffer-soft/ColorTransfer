package tmf

import (
	"archive/zip"
	"bytes"
)

func Pack3mf(modelXml, slic3rConfig, slic3rModelConfig string) []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	addFile := func(name, data string) {
		f, _ := zipWriter.Create(name)
		_, _ = f.Write([]byte(data))
	}

	addFile("_rels/.rels", RelsXml)
	addFile("3D/3dmodel.model", modelXml)
	addFile("Metadata/Slic3r_PE.config", slic3rConfig)
	addFile("Metadata/Slic3r_PE_model.config", slic3rModelConfig)
	addFile("[Content_Types].xml", ContentTypesXml)

	err := zipWriter.Close()
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
