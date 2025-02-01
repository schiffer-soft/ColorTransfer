package main

import (
	"colortransfer_cmd/tmf"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	//os.Args = append(os.Args, "Crew 9 v101.3mf")
	if len(os.Args) < 2 {
		fmt.Println("Usage: colortransfer <inputfile> <outputfile> hexcodes")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println()
		fmt.Println("  --- display colors ---")
		fmt.Println("  colortransfer test.3mf")
		fmt.Println()
		fmt.Println("  --- convert to default extruders ---")
		fmt.Println("  colortransfer test.3mf newtest.3mf")
		fmt.Println()
		fmt.Println("  --- convert to color-sorted extruders ---")
		fmt.Println("  colortransfer test.3mf newtest.3mf #002BA2;#092A60;#5375C3;#EEE13E;#87A9EC")
		fmt.Println()
		os.Exit(1)
		return
	}
	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: File not found:", os.Args[1])
		os.Exit(1)
		return
	}
	model, err := tmf.GetModelFrom3mf(bytes)
	if err != nil {
		fmt.Println("Error: Decode Error:", err)
		os.Exit(1)
		return
	}
	if len(model.Resources.ColorGroups) < 2 {
		fmt.Println("Error: No 3MF-Colors found.")
		os.Exit(1)
		return
	}

	var modelXml, slic3rConfig, slic3rModelConfig string

	var verts []tmf.Vertex
	var tris []tmf.Triangle
	vertOfs := 0
	config := tmf.Config{Objects: []tmf.ConfigObject{{ID: 1, InstancesCount: 1, Metadata: []tmf.Metadata2{{Type: "object", Key: "name", Value: model.Metadata[0].Value}}}}}

	extruders := map[string]int{}

	for _, color := range model.Resources.ColorGroups {
		c := color.Colors[0].Color[:7]
		extruder := extruders[c]
		if extruder == 0 {
			extruders[c] = len(extruders) + 1
			extruder = extruders[c]
			fmt.Println(c)
		}
	}
	fmt.Println()

	if len(os.Args) < 3 {
		png := tmf.GetThumbnailFrom3mf(bytes)
		fmt.Println(base64.StdEncoding.EncodeToString(png))
		fmt.Println()
		return
	}

	if len(os.Args) > 3 {
		sp := strings.Split(os.Args[3], ";")
		if len(sp) > 0 && sp[len(sp)-1] == "" {
			sp = sp[:len(sp)-1]
		}
		known := make(map[string]bool)
		for i, check := range sp {
			if known[check] {
				fmt.Println("Error: Duplicate Extruder-Color found:", check)
				os.Exit(1)
				return
			}
			known[check] = true
			if extruders[check] == 0 {
				fmt.Println("Error: Extruder-Color not found:", check)
				os.Exit(1)
				return
			}
			extruders[check] = i + 1
		}
		if len(sp) != len(extruders) {
			fmt.Println("Error: Extruder-Count not match:", len(sp), "vs", len(extruders))
			os.Exit(1)
			return
		}
		fmt.Println("--- sorted extruders ---")
		for i := range sp {
			fmt.Println(sp[i], "-", extruders[sp[i]])
		}
		fmt.Println("---")
	}

	if err = os.WriteFile(os.Args[2], []byte("hello"), 0644); err != nil {
		fmt.Println("Error: WriteFile Error:", err)
		fmt.Println(os.Args[2])
		os.Exit(1)
		return
	}

	for objIndex, obj := range model.Resources.Objects {
		extruder := extruders[model.Resources.ColorGroups[objIndex].Colors[0].Color[:7]]
		verts = append(verts, obj.Mesh.Vertices...)
		for _, tri := range obj.Mesh.Triangles {
			tris = append(tris, tmf.Triangle{V1: vertOfs + tri.V1, V2: vertOfs + tri.V2, V3: vertOfs + tri.V3})
		}
		config.Objects[0].Volumes = append(config.Objects[0].Volumes, tmf.Volume{
			FirstID:  len(tris) - len(obj.Mesh.Triangles),
			LastID:   len(tris) - 1,
			Metadata: []tmf.Metadata2{{"volume", "name", obj.Name}, {"volume", "volume_type", "ModelPart"}, {"volume", "matrix", "1 0 0 0 0 1 0 0 0 0 1 0 0 0 0 1"}, {"volume", "extruder", strconv.Itoa(extruder)}},
			Mesh:     tmf.MeshInfo{},
		})
		vertOfs = len(verts)
	}
	var model2 tmf.Model
	model2.XMLName = model.XMLName
	model2.Metadata = model.Metadata
	model2.Resources.Objects = append(model2.Resources.Objects, tmf.Object{ID: 1, Type: "model", Mesh: tmf.Mesh{Vertices: verts, Triangles: tris}})

	if dest, err := xml.Marshal(model2); err != nil {
		log.Fatal(err)
	} else {
		name := model2.Metadata[0].Value
		modelXml = string(dest)[:len(dest)-8] + "<build><item objectid=\"1\" transform=\"1 0 0 0 1 0 0 0 1 175 180 60\" printable=\"1\"/></build></model>"
		modelXml = `<?xml version="1.0" encoding="UTF-8"?>
	<model unit="millimeter" xml:lang="en-US" xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:slic3rpe="http://schemas.slic3r.org/3mf/2017/06">
	<metadata name="slic3rpe:Version3mf">1</metadata>
	<metadata name="Title">` + name + `</metadata>
	<metadata name="Designer"></metadata>
	<metadata name="Description">` + name + `</metadata>
	<metadata name="Copyright"></metadata>
	<metadata name="LicenseTerms"></metadata>
	<metadata name="Rating"></metadata>
	<metadata name="CreationDate">2024-10-09</metadata>
	<metadata name="ModificationDate">2024-10-09</metadata>
	<metadata name="Application">PrusaSlicer-2.8.1+win64</metadata>
	` + modelXml[7:]
		modelXml = strings.ReplaceAll(modelXml, "></vertex>", " />")
		modelXml = strings.ReplaceAll(modelXml, "></triangle>", " />")
		fmt.Println("3dmodel.model", len(modelXml), "bytes")
	}
	if dest, err := xml.Marshal(config); err != nil {
		log.Fatal(err)
	} else {
		slic3rModelConfig = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" + string(dest)
		slic3rModelConfig = strings.ReplaceAll(slic3rModelConfig, "></mesh>", " />")
		slic3rModelConfig = strings.ReplaceAll(slic3rModelConfig, "></metadata>", " />")
		//_ = os.WriteFile("Slic3r_PE_model.config", []byte(str), 0644)
		fmt.Println("Slic3r_PE_model.config", len(slic3rModelConfig), "bytes")
	}
	{
		slic3rConfig = ""
		for i := 1; i <= len(extruders); i++ {
			for k, v := range extruders {
				if v == i {
					slic3rConfig += k + ";"
				}
			}
		}
		for i := len(extruders); i < 5; i++ {
			slic3rConfig += "#000000;"
		}
		slic3rConfig = slic3rConfig[:len(slic3rConfig)-1]
		fmt.Println("extruder:", slic3rConfig)
		slic3rConfig = strings.ReplaceAll(tmf.BaseCfg, "{{colors}}", slic3rConfig)
		fmt.Println("Slic3r_PE.config", len(slic3rConfig), "bytes")
	}
	fmt.Println("packe 3mf...")
	output := tmf.Pack3mf(modelXml, slic3rConfig, slic3rModelConfig)
	fmt.Println("packe 3mf: ok", len(output), "bytes")

	if err = os.WriteFile(os.Args[2], output, 0644); err != nil {
		fmt.Println("Error: WriteFile Error:", err)
		fmt.Println(os.Args[2])
		os.Exit(1)
		return
	}

	fmt.Println("ok.")
}
