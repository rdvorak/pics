package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

type Options struct {
	monthName map[string]string
	link      string
}
type Word struct {
	Text   string `json:"text"`
	Weight int    `json:"weight"`
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
	options.link = "http://localhost:8081/gallery"
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
	router.StaticFS("/jqcloud", http.Dir("jqcloud"))
	router.StaticFS("/jquery", http.Dir("jquery"))
	router.StaticFS("/folio", http.Dir("folio"))
	router.StaticFS("/pics", http.Dir("/Users/rdvorak/Pictures"))

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
	router.GET("/gallery/drilldown", func(c *gin.Context) {
		var tags []interface{}
		fmt.Println(c.QueryArray("tag"))
		for _, v := range c.QueryArray("tag") {
			fmt.Println(v)
			tags = append(tags, v)
		}
		c.HTML(http.StatusOK, "index.html", db.drillByTags(tags...))

		// b, _ := json.Marshal(db.drillByTags(tags...).tags)
		// c.String(200, "%s", b)
	})
	router.Run(":8081")
}
