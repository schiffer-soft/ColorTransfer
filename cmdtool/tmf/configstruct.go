package tmf

import "encoding/xml"

// Config represents the root element of the XML structure
type Config struct {
	XMLName xml.Name       `xml:"config"`
	Objects []ConfigObject `xml:"object"`
}

type ConfigObject struct {
	ID             int         `xml:"id,attr"`
	InstancesCount int         `xml:"instances_count,attr"`
	Metadata       []Metadata2 `xml:"metadata"`
	Volumes        []Volume    `xml:"volume"`
}

type Volume struct {
	FirstID  int         `xml:"firstid,attr"`
	LastID   int         `xml:"lastid,attr"`
	Metadata []Metadata2 `xml:"metadata"`
	Mesh     MeshInfo    `xml:"mesh"`
}

type Metadata2 struct {
	Type  string `xml:"type,attr"`
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

type MeshInfo struct {
	EdgesFixed       int `xml:"edges_fixed,attr"`
	DegenerateFacets int `xml:"degenerate_facets,attr"`
	FacetsRemoved    int `xml:"facets_removed,attr"`
	FacetsReversed   int `xml:"facets_reversed,attr"`
	BackwardsEdges   int `xml:"backwards_edges,attr"`
}
