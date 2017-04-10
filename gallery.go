package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
)

type Options struct {
	monthName       map[string]string
	link            string
	Source          string
	CloudTags       []string
	descriptionTags []string
}
type Word struct {
	Text   string `json:"text"`
	Count  int    `json:"weight"`
	Weight int    `json:"weightgrp"`
	Link   string `json:"link"`
}
type Image struct {
	Image       string `json:"image,omitempty"`
	Thumb       string `json:"thumb,omitempty"`
	Big         string `json:"big,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
}
type Gallery struct {
	Tags   []Word
	Images []Image
}

var options Options

func main() {
	options.monthName = map[string]string{"01": "Leden", "02": "Únor", "03": "Březen", "04": "Duben", "05": "Květen", "06": "Červen", "07": "Červenec", "08": "Srpen", "09": "Září", "10": "Říjen", "11": "Listopad", "12": "Prosinec"}
	options.link = "/gallery"
	options.Source = "/Users/rdvorak/Pictures"
	options.CloudTags = []string{"Rating", "FocalLength", "Keyword", "Month", "Year", "State", "Country", "Location", "Sublocation", "Geoname"}
	options.descriptionTags = []string{"Rating", "Model", "Lens", "Keyword", "Month", "Year", "State", "Country", "Location", "Sublocation", "Geoname"}
	db := pictureDb()
	defer db.sess.Close()
	// allGallery := db.drillByTags()
	// fmt.Printf("%v\n", allGallery)
	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	gin.SetMode("debug")
	router := gin.Default()
	// admin := router.Group("/admin/", gin.BasicAuth(gin.Accounts{
	// "pozdechov": "vp2",
	// }))
	// router.StaticFS("/jqcloud", http.Dir("jqcloud"))
	router.StaticFS("/jquery", http.Dir("jquery"))
	router.StaticFS("/folio", http.Dir("folio"))
	router.StaticFS("/pics", http.Dir(options.Source))

	// router.StaticFS("/Bacovi-rodokmen", http.Dir("www/Bacovi-rodokmen"))
	// router.LoadHTMLFiles("www/vp2.html")
	// router.GET("/", func(c *gin.Context) {
	// c.HTML(http.StatusOK, "vp2.html", nil)
	// })

	// router.GET("/archive/day", func(c *gin.Context) {
	// c.JSON(http.StatusOK, vp.archiveDay)
	// })
	router.POST("/submit/metadata", func(c *gin.Context) {
		var data []Metadata
		err := c.BindJSON(&data)
		if err != nil {
			log.Println(err)
		}
		for _, meta := range data {
			p := Picture{}
			p.ParseMetadata(meta)
			if meta.GPSLatitude != "" && meta.GPSLongitude != "" {
				savedLocation := db.getLocation(meta.GPSLatitude, meta.GPSLongitude)
				var location *OsmAddress
				if savedLocation == "" {
					location = OsmLocation(meta.GPSLatitude, meta.GPSLongitude)
					time.Sleep(time.Millisecond * 500)
				} else {
					err = json.Unmarshal([]byte(savedLocation), location)
					if err != nil {
						fmt.Printf("OSM address: unmarshall error: %v\n %v\n", err, savedLocation)
					}
				}
				p.ParseLocation(location)

			}
			db.savePicture(p)
		}
	})
	router.POST("/submit/metadata/num", func(c *gin.Context) {
		var data []NumMetadata
		err := c.BindJSON(&data)
		if err != nil {
			log.Println(err)
		}
		for _, meta := range data {
			p := Picture{}
			p.ParseNumMetadata(meta)
			db.savePicture(p)
		}
	})
	router.LoadHTMLFiles("index.html")
	router.GET("/gallery", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", db.drillByTags(""))
	})
	router.GET("/gallery/query", func(c *gin.Context) {
		var tags []interface{}
		for _, v := range c.QueryArray("tag") {
			// fmt.Println(v)
			tags = append(tags, v)
		}
		c.HTML(http.StatusOK, "index.html", db.drillByTags(tags...))

	})
	router.Run(":8081")
}
