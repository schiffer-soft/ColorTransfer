package tmf

import "encoding/xml"

type Model struct {
	XMLName   xml.Name   `xml:"model"`
	Metadata  []Metadata `xml:"metadata"`
	Resources Resources  `xml:"resources"`
}

type Metadata struct {
	Type     string `xml:"type,attr"`
	Key      string `xml:"key,attr"`
	Preserve string `xml:"preserve,attr"`
	Value    string `xml:",chardata"`
}

type Resources struct {
	ColorGroups []ColorGroup `xml:"colorgroup"`
	Objects     []Object     `xml:"object"`
}

type ColorGroup struct {
	ID     int     `xml:"id,attr"`
	Colors []Color `xml:"color"`
}

type Color struct {
	Color string `xml:"color,attr"`
}

type Object struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Mesh Mesh   `xml:"mesh"`
}

type Mesh struct {
	Vertices  []Vertex   `xml:"vertices>vertex"`
	Triangles []Triangle `xml:"triangles>triangle"`
}

type Vertex struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
	Z float64 `xml:"z,attr"`
}

type Triangle struct {
	V1 int `xml:"v1,attr"`
	V2 int `xml:"v2,attr"`
	V3 int `xml:"v3,attr"`
}
