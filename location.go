package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type mapyczRgeocode struct {
	XMLName xml.Name `xml:"rgeocode"`
	Items   []struct {
		Name     string `xml:"name,attr"`
		ItemType string `xml:"type,attr"`
	} `xml:"item"`
}
type OsmAddress struct {
	Address struct {
		HouseNumber string `json:"house_number"`
		Building    string `json:"building"`
		Path        string `json:"path"`
		Road        string `json:"road"`
		Suburb      string `json:"suburb"`
		Village     string `json:"village"`
		City        string `json:"city"`
		Town        string `json:"town"`
		County      string `json:"county"`
		State       string `json:"state"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
	} `json:"address"`
	NameDetails struct {
		Name string `json:"name"`
	} `json:"namedetails"`
}

func mapyczLocation(lat, lon string) map[string]string {
	if strings.HasSuffix(lon, " W") {
		lat = "-" + strings.TrimSuffix(lon, " W")
	}
	if strings.HasSuffix(lat, " S") {
		lat = "-" + strings.TrimSuffix(lat, " S")
	}
	resp, err := http.Get("https://api.mapy.cz/rgeocode?" + "lat=" + lat + "&lon=" + lon)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}
	fmt.Printf("%s", body)
	rgeocode := mapyczRgeocode{}

	err = xml.Unmarshal(body, &rgeocode)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}
	loc := make(map[string]string)
	for _, item := range rgeocode.Items {
		switch item.ItemType {
		case "ward":
			loc["Sublocation"] = item.Name
		case "muni":
			loc["Location"] = item.Name
			if loc["Sublocation"] == loc["Location"] {
				delete(loc, "Sublocation")
			}
		case "regi":
			loc["Region"] = item.Name
		case "coun":
			loc["Country"] = item.Name
		}
	}
	return loc

}
func OsmLocation(lat, lon string) *OsmAddress {
	if strings.HasSuffix(lat, " N") {
		lat = strings.TrimSuffix(lat, " N")
	}
	if strings.HasSuffix(lat, " S") {
		lat = "-" + strings.TrimSuffix(lat, " S")
	}
	if strings.HasSuffix(lon, " E") {
		lon = strings.TrimSuffix(lon, " E")
	}
	if strings.HasSuffix(lon, " W") {
		lon = "-" + strings.TrimSuffix(lon, " W")
	}
	resp, err := http.Get("http://nominatim.openstreetmap.org/reverse?format=json&namedetails=1&zoom=18&" + "lat=" + lat + "&lon=" + lon)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	// fmt.Printf("%s", body)
	address := OsmAddress{}

	err = json.Unmarshal(body, &address)
	if err != nil {
		fmt.Printf("OSM address: unmarshall error: %v\n %s\n", err, body)
	}
	return &address

}
