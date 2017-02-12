package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/gin-gonic/gin.v1"
)

type Options struct {
	monthName map[string]string
}
type Gallery struct {
	tags     map[string]int
	pictures map[string][]byte
}

var options Options

func main() {
	options.monthName = map[string]string{"01": "Leden", "02": "Únor", "03": "Březen", "04": "Duben", "05": "Květen", "06": "Červen", "07": "Červenec", "08": "Srpen", "09": "Září", "10": "Říjen", "11": "Listopad", "12": "Prosinec"}
	db := pictureDb()
	defer db.sess.Close()
	allGallery := db.drillByTags()
	fmt.Printf("%v\n", allGallery)
	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	gin.SetMode("debug")
	router := gin.Default()
	// admin := router.Group("/admin/", gin.BasicAuth(gin.Accounts{
	// "pozdechov": "vp2",
	// }))
	// router.StaticFS("/static", http.Dir("www"))
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
	router.GET("/gallery/drilldown", func(c *gin.Context) {
		var tags []interface{}
		fmt.Println(c.Params)
		for _, v := range c.QueryArray("tag") {
			fmt.Println(v)
			tags = append(tags, v)
		}
		b, _ := json.Marshal(db.drillByTags(tags...).tags)
		c.String(200, "%s", b)
	})
	router.Run(":8081")
}
