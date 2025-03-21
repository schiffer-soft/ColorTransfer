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

func StrFl(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

func Println(a ...any) {
	o := fmt.Appendln(nil, a...)
	o = append(o[:len(o)-1], '\r', '\n')
	fmt.Print(string(o))
}

func main() {
	//os.Args = append(os.Args, "Crew9.3mf", "Crew9_mk4s2.3mf", "Form-Kubus.3mf", "#002BA2;#092A60;#5375C3;#87A9EC;#EEE13E")
	//os.Args = append(os.Args, "Crew9.3mf", "Crew9_fix.3mf")
	if len(os.Args) < 2 {
		Println("Usage: colortransfer <inputfile> <outputfile> hexcodes")
		Println()
		Println("Examples:")
		Println()
		Println("  --- display colors ---")
		Println("  colortransfer test.3mf")
		Println()
		Println("  --- convert to default extruders ---")
		Println("  colortransfer test.3mf newtest.3mf")
		Println()
		Println("  --- convert to default extruders and source profile ---")
		Println("  colortransfer test.3mf newtest.3mf profile.3mf")
		Println()
		Println("  --- convert to color-sorted extruders ---")
		Println("  colortransfer test.3mf newtest.3mf #002BA2;#092A60;#5375C3;#87A9EC;#EEE13E")
		Println()
		Println("  --- convert to color-sorted extruders and source profile ---")
		Println("  colortransfer test.3mf newtest.3mf profile.3mf #002BA2;#092A60;#5375C3;#87A9EC;#EEE13E")
		Println()
		os.Exit(1)
		return
	}
	if len(os.Args) == 4 && !strings.HasSuffix(os.Args[len(os.Args)-1], ".3mf") {
		os.Args = append(os.Args[:len(os.Args)-1], "", os.Args[len(os.Args)-1])
	}
	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		Println("Error: File not found:", os.Args[1])
		os.Exit(1)
		return
	}
	model, err := tmf.GetModelFrom3mf(bytes)
	if err != nil {
		Println("Error: Decode Error:", err)
		os.Exit(1)
		return
	}
	if len(model.Resources.ColorGroups) < 1 {
		Println("Error: No 3MF-Colors found.")
		os.Exit(1)
		return
	}

	for i := len(model.Resources.Objects) - 1; i >= 0; i-- {
		if len(model.Resources.Objects[i].Mesh.Triangles) == 0 {
			model.Resources.ColorGroups = append(model.Resources.ColorGroups[:i], model.Resources.ColorGroups[i+1:]...)
			model.Resources.Objects = append(model.Resources.Objects[:i], model.Resources.Objects[i+1:]...)
		}
	}
	if len(model.Resources.ColorGroups) < 1 {
		Println("Error: No Models found.")
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
			Println(c)
		}
	}
	Println()

	if len(os.Args) < 3 {
		png := tmf.GetThumbnailFrom3mf(bytes)
		Println("picture:", base64.StdEncoding.EncodeToString(png))
		Println()
		return
	}

	if len(os.Args) > 4 {
		sp := strings.Split(os.Args[4], ";")
		if len(sp) > 0 && sp[len(sp)-1] == "" {
			sp = sp[:len(sp)-1]
		}
		known := make(map[string]bool)
		for i, check := range sp {
			if known[check] {
				Println("Error: Duplicate Extruder-Color found:", check)
				os.Exit(1)
				return
			}
			if extruders[check] == 0 {
				if len(known) == len(extruders) {
					sp = sp[:len(known)] // trunc too many colors
					break
				}
				Println("Error: Extruder-Color not found:", check)
				os.Exit(1)
				return
			}
			known[check] = true
			extruders[check] = i + 1
		}
		if len(sp) != len(extruders) {
			Println("Error: Extruder-Count not match:", len(sp), "vs", len(extruders))
			os.Exit(1)
			return
		}
		Println("--- sorted extruders ---")
		for i := range sp {
			Println(sp[i], "-", extruders[sp[i]])
		}
		Println("---")
	}

	if err = os.WriteFile(os.Args[2], []byte("hello"), 0644); err != nil {
		Println("Error: WriteFile Error:", err)
		Println(os.Args[2])
		os.Exit(1)
		return
	}

	for objIndex, obj := range model.Resources.Objects {
		extruder := extruders[model.Resources.ColorGroups[objIndex].Colors[0].Color[:7]]
		verts = append(verts, obj.Mesh.Vertices...)
		for _, tri := range obj.Mesh.Triangles {
			tris = append(tris, tmf.Triangle{V1: vertOfs + tri.V1, V2: vertOfs + tri.V2, V3: vertOfs + tri.V3})
		}
		if len(model.Resources.ColorGroups[objIndex].Colors) > 2 {
			Println("too many colors per object:", obj.Name, "-", len(model.Resources.ColorGroups[objIndex].Colors)-1)
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

	// --- read profile ---
	slic3rProfile := tmf.BaseCfg
	bedWidth := 360.0  // XL Default
	bedHeight := 360.0 // XL Default
	const bedShape = "; bed_shape = "
	const extColor = "; extruder_colour = "
	const filColor = "; filament_colour = "
	if len(os.Args) > 3 {
		if tmp := tmf.GetModelFrom3mfFile(os.Args[3]); strings.Contains(tmp, bedShape) && strings.Contains(tmp, extColor) && strings.Contains(tmp, filColor) {
			p := strings.Index(tmp, bedShape)
			p2 := strings.IndexByte(tmp[p:], '\n') + p
			if sp := strings.Split(tmp[p+len(bedShape):p2], ","); len(sp) == 4 {
				sp = strings.Split(sp[2], "x")
				if len(sp) == 2 {
					w, _ := strconv.ParseFloat(sp[0], 64)
					h, _ := strconv.ParseFloat(sp[1], 64)
					if w > 10 && w < 100000 {
						bedWidth = w
					}
					if h > 10 && h < 100000 {
						bedHeight = h
					}
				}
			}
			p = strings.Index(tmp, extColor)
			p2 = strings.IndexByte(tmp[p:], '\n') + p
			tmp = tmp[:p] + extColor + "{{colors}}" + tmp[p2:]
			p = strings.Index(tmp, filColor)
			p2 = strings.IndexByte(tmp[p:], '\n') + p
			tmp = tmp[:p] + filColor + "{{colors}}" + tmp[p2:]
			slic3rProfile = tmp
		}
	}

	centerX := 0.0
	centerY := 0.0
	lowZ := 1000000.0
	for _, v := range verts {
		centerX += v.X
		centerY += v.Y
		if v.Z < lowZ {
			lowZ = v.Z
		}
	}
	centerX /= float64(len(verts))
	centerY /= float64(len(verts))

	if dest, err := xml.Marshal(model2); err != nil {
		log.Fatal(err)
	} else {
		name := model2.Metadata[0].Value
		modelXml = string(dest)[:len(dest)-8] + "<build><item objectid=\"1\" transform=\"1 0 0 0 1 0 0 0 1 " + StrFl(bedWidth*0.5-centerX) + " " + StrFl(bedHeight*0.5-centerY) + " " + StrFl(-lowZ) + "\" printable=\"1\"/></build></model>"
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
	<metadata name="Application">PrusaSlicer-2.9.1+win64</metadata>
	` + modelXml[7:]
		modelXml = strings.ReplaceAll(modelXml, "></vertex>", " />")
		modelXml = strings.ReplaceAll(modelXml, "></triangle>", " />")
		Println("3dmodel.model", len(modelXml), "bytes")
	}
	if dest, err := xml.Marshal(config); err != nil {
		log.Fatal(err)
	} else {
		slic3rModelConfig = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" + string(dest)
		slic3rModelConfig = strings.ReplaceAll(slic3rModelConfig, "></mesh>", " />")
		slic3rModelConfig = strings.ReplaceAll(slic3rModelConfig, "></metadata>", " />")
		//_ = os.WriteFile("Slic3r_PE_model.config", []byte(str), 0644)
		Println("Slic3r_PE_model.config", len(slic3rModelConfig), "bytes")
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
		Println("extruder:", slic3rConfig)
		slic3rConfig = strings.ReplaceAll(slic3rProfile, "{{colors}}", slic3rConfig)
		Println("Bed-Size:", bedWidth, "x", bedHeight)
		Println("Slic3r_PE.config", len(slic3rConfig), "bytes")
	}
	Println("packe 3mf...")
	output := tmf.Pack3mf(modelXml, slic3rConfig, slic3rModelConfig)
	Println("packe 3mf: ok", len(output), "bytes")

	if err = os.WriteFile(os.Args[2], output, 0644); err != nil {
		Println("Error: WriteFile Error:", err)
		Println(os.Args[2])
		os.Exit(1)
		return
	}

	Println("ok.")
}
