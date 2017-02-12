package main

import (
	"strings"
)

type Metadata struct {
	FileName             string
	Model                string
	DateTimeOriginal     string
	ISO                  int
	ShutterSpeed         string
	Aperture             float64
	ExposureCompensation float64
	FocalLength          string
	WhiteBalance         string
	ImageSize            string
	Orientation          string
	LensInfo             string
	LensID               string
	Lens                 string
	GPSAltitude          string
	GPSLatitude          string
	GPSLongitude         string
	Rating               int
	Keywords             interface{}
}
type NumMetadata struct {
	FileName             string
	DateTimeOriginal     string
	ISO                  int
	ShutterSpeed         float64
	Aperture             float64
	ExposureCompensation float64
	FocalLength          float64
	WhiteBalance         float64
	Orientation          int
	GPSAltitude          float64
	GPSLatitude          float64
	GPSLongitude         float64
	Rating               int
}

type Tag struct {
	lang   string
	tag    string
	meta   string
	source int
}

type Picture struct {
	Name  string
	Year  string
	Month string
	*Metadata
	*NumMetadata
	Tags []Tag
}

func (p *Picture) ParseMetadata(m Metadata) {
	t := strings.Split(strings.Split(m.DateTimeOriginal, " ")[0], ":")
	p.Year = t[0]

	p.Month = t[1]
	p.Name = t[0] + "/" + t[1] + "/" + m.FileName
	var tags []Tag
	tags = append(tags, Tag{meta: "Model", tag: m.Model})
	tags = append(tags, Tag{meta: "FocalLength", tag: m.FocalLength})
	// tags = append(tags, Tag{meta: "ShutterSpeed", tag: m.ShutterSpeed})
	// tags = append(tags, Tag{meta: "Aperture", tag: m.Aperture})
	if m.Lens == "" {
		if m.FocalLength == "35.0 mm" && strings.HasPrefix(m.Model, "NIKON") {

			tags = append(tags, Tag{meta: "Lens", tag: "Zeiss 35 mm"})
		}
	} else {
		tags = append(tags, Tag{meta: "Lens", tag: m.Lens})
	}
	tags = append(tags, Tag{meta: "Rating", tag: strings.Repeat("*", m.Rating)})
	switch m.Keywords.(type) {
	case []string:
		for _, k := range m.Keywords.([]string) {
			tags = append(tags, Tag{meta: "Keyword", tag: k})
		}
	case string:

		tags = append(tags, Tag{meta: "Keyword", tag: m.Keywords.(string)})
	}
	tags = append(tags, Tag{meta: "Year", tag: p.Year})
	tags = append(tags, Tag{meta: "Month", tag: options.monthName[p.Month]})
	for i := range tags {
		tags[i].source = 1
	}
	p.Tags = tags
	p.Metadata = &m
}
func (p *Picture) ParseNumMetadata(m NumMetadata) {
	if p.Name == "" {
		t := strings.Split(strings.Split(m.DateTimeOriginal, " ")[0], ":")
		p.Name = t[0] + "/" + t[1] + "/" + t[2] + "/" + m.FileName
	}
	p.NumMetadata = &m
}
