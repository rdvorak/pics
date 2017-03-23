package main

import (
	"fmt"
	"log"
	"strings"
	"time"
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
	Location *OsmAddress
	Tags     []Tag
}

func (p *Picture) ParseMetadata(m Metadata) {
	t := strings.Split(strings.Split(m.DateTimeOriginal, " ")[0], ":")
	if len(t) > 1 {

		p.Year = t[0]
		p.Month = t[1]
		p.Name = t[0] + "/" + t[1] + "/" + m.FileName
	} else {
		log.Println("picture metadata: DateTimeOriginal: ", m.DateTimeOriginal, " cannot be parsed")
		return
	}
	var tags []Tag
	tags = append(tags, Tag{meta: "Model", tag: m.Model})
	tags = append(tags, Tag{meta: "FocalLength", tag: m.FocalLength})
	// tags = append(tags, Tag{meta: "ShutterSpeed", tag: m.ShutterSpeed})
	// tags = append(tags, Tag{meta: "Aperture", tag: m.Aperture})

	if m.Lens == "" {
		if m.FocalLength == "35.0 mm" && strings.HasPrefix(m.Model, "NIKON") {

			tags = append(tags, Tag{meta: "Lens", tag: "Zeiss Distagon 35mm f/2 T* ZF"})
		}
	} else {
		tags = append(tags, Tag{meta: "Lens", tag: m.Lens})
	}

	tags = append(tags, Tag{meta: "Rating", tag: strings.Repeat("*", m.Rating) + " Star"})
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
	if m.GPSLatitude != "" && m.GPSLongitude != "" {
		m.GPSLatitude = strings.TrimPrefix(m.GPSLatitude, "+")
		m.GPSLongitude = strings.TrimPrefix(m.GPSLongitude, "+")
		addr := OsmLocation(m.GPSLatitude, m.GPSLongitude)

		if addr.Address.Suburb != "" {
			tags = append(tags, Tag{meta: "Sublocation", tag: addr.Address.City + addr.Address.Village, source: 3})
		}
		if addr.Address.City != "" || addr.Address.Village != "" {
			tags = append(tags, Tag{meta: "Location", tag: addr.Address.City + addr.Address.Village, source: 3})
		}
		if addr.Address.State != "" {
			tags = append(tags, Tag{meta: "State", tag: addr.Address.State, source: 3})
		}
		if addr.NameDetails.Name != "" {
			tags = append(tags, Tag{meta: "Geoname", tag: addr.NameDetails.Name, source: 3})
		}
		if addr.Address.Country != "" {
			tags = append(tags, Tag{meta: "Country", tag: addr.Address.Country, source: 3})
		}
		time.Sleep(time.Millisecond * 500)
		p.Location = &addr
	}
	p.Tags = tags
	p.Metadata = &m
	fmt.Printf("%#v\n", p)
}
func (p *Picture) ParseNumMetadata(m NumMetadata) {
	if p.Name == "" {
		t := strings.Split(strings.Split(m.DateTimeOriginal, " ")[0], ":")
		p.Name = t[0] + "/" + t[1] + "/" + t[2] + "/" + m.FileName
	}
	p.NumMetadata = &m
}
