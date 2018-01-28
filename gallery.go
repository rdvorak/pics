package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
)

type Word struct {
	Text   string `json:"text"`
	Count  int    `json:"weight"`
	Weight int    `json:"weightgrp"`
	Link   string `json:"link"`
	Color  string `json:"color"`
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

func main() {
	db := pictureDb()
	defer db.sess.Close()
	gin.SetMode("debug")
	router := gin.Default()
	router.StaticFS("/gallery/jquery", http.Dir("jquery"))
	router.StaticFS("/gallery/lightgallery", http.Dir("lightgallery"))
	router.StaticFS("/gallery/folio", http.Dir("folio"))
	router.StaticFS("/pics", http.Dir(options.Pics))
	router.StaticFS("/thums", http.Dir(options.Thums))

	router.POST("/gallery/submit/metadata", func(c *gin.Context) {
		var data []Metadata
		err := c.BindJSON(&data)
		if err != nil {
			log.Println(err)
		}
		for _, meta := range data {
			p := Picture{}
			p.ParseMetadata(&meta)
			if meta.GPSLatitude != "" && meta.GPSLongitude != "" {
				savedLocation := db.getLocation(meta.GPSLatitude, meta.GPSLongitude)
				var location *OsmAddress
				if savedLocation == "" {
					fmt.Printf("getting OSM adress: %s lat=%s,lon=%s \n", p.Name, meta.GPSLatitude, meta.GPSLongitude)
					location = OsmLocation(meta.GPSLatitude, meta.GPSLongitude)
					time.Sleep(time.Millisecond * 500)
				} else {
					err = json.Unmarshal([]byte(savedLocation), &location)
					if err != nil {
						fmt.Printf("OSM address: unmarshall error: %v\n %s\n", err, savedLocation)
					}
				}
				p.ParseLocation(location)

			}
			db.savePicture(p)
		}
	})
	router.POST("/gallery/submit/metadata/num", func(c *gin.Context) {
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
	gallery := db.drillByTags()
	randImages := func(im []Image) []Image {
		var sel []Image
		// nahodne vybereme 100 obrazku
		// fmt.Println("Images ", len(sel.Images))
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i, val := range r.Perm(len(im)) {
			sel = append(sel, im[val])
			if i > 100 {
				break
			}
		}
		return sel
	}
	router.GET("/gallery", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", Gallery{Tags: gallery.Tags, Images: randImages(gallery.Images)})
	})
	router.GET("/gallery/query", func(c *gin.Context) {
		var tags []interface{}
		for _, v := range c.QueryArray("tag") {
			// fmt.Println(v)
			tags = append(tags, v)
		}

		sel := db.drillByTags(tags...)
		c.HTML(http.StatusOK, "index.html", Gallery{Tags: sel.Tags, Images: randImages(sel.Images)})

	})
	router.Run(":8081")
}
