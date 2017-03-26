package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type PictureDb struct {
	sess *sql.DB
}

func pictureDb() PictureDb {
	db, err := sql.Open("sqlite3", "./pics.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
create table if not exists Pictures (
	NAME text not null primary key,
	FILENAME text,
	METADATA	text,
	NUM_METADATA  text,
	LOCATION text
);
create table if not exists Picture_Tags (
	picture_name text not null ,
	meta text, tag text, lang text, source int
);
create unique index if not exists picture_tags_x1 on picture_tags (picture_name, tag);
create unique index if not exists picture_tags_x2 on picture_tags (tag, picture_name);
`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Panicf("%q: %s\n", err, sqlStmt)
	}
	return PictureDb{db}
}
func (db *PictureDb) savePicture(p Picture) {
	if p.Metadata != nil {
		stmt := `insert or replace into Pictures ( NAME, FILENAME, METADATA ) values( ?, ?, ?)`
		meta, _ := json.Marshal(p.Metadata)
		_, err := db.sess.Exec(stmt, p.Name, p.Metadata.FileName, meta)
		if err != nil {
			log.Panicln(err)
		}
		_, err = db.sess.Exec(`delete from PICTURE_TAGS where picture_name = ? and source = 1`, p.Name)
		if err != nil {
			log.Panicln(err)
		}
	}
	if p.NumMetadata != nil {
		numMeta, _ := json.Marshal(p.NumMetadata)
		_, err := db.sess.Exec(`update pictures set num_metadata = ? where name = ?`, numMeta, p.Name)
		if err != nil {
			log.Panicln(err)
		}
		_, err = db.sess.Exec(`delete from PICTURE_TAGS where picture_name = ? and source = 2`, p.Name)
		if err != nil {
			log.Panicln(err)
		}
	}

	if p.Location != nil {
		loc, _ := json.Marshal(p.Location)
		_, err := db.sess.Exec(`update pictures set location = ? where name = ?`, loc, p.Name)
		if err != nil {
			log.Panicln(err)
		}
		_, err = db.sess.Exec(`delete from PICTURE_TAGS where picture_name = ? and (source = 3 or source is null)`, p.Name)
		if err != nil {
			log.Panicln(err)
		}
	}
	if len(p.Tags) > 0 {
		stmt := `insert or replace into Picture_tags ( PICTURE_NAME, META, TAG, LANG, SOURCE ) values( ?, ?, ?, ?, ? )`
		for _, tag := range p.Tags {
			if tag.tag != "" {
				_, err := db.sess.Exec(stmt, p.Name, tag.meta, tag.tag, tag.lang, tag.source)
				if err != nil {
					log.Panicln(err)
				}
			}
		}
	}
}
func (db *PictureDb) drillByTags(tags ...interface{}) Gallery {
	var sel Gallery
	var sql string
	var params string
	var CloudTags string
	if len(options.CloudTags) > 0 {
		CloudTags = " and meta in ('" + strings.Join(options.CloudTags, "','") + "')"
	}
	for i, tag := range tags {
		if i == 0 {
			sql = "select picture_name from picture_tags where tag = ?  " + CloudTags
		} else if i < len(tags) {
			sql = "select picture_name from picture_tags where picture_name in (" + sql + ") and tag = ? " + CloudTags
		}
		params = params + "&tag=" + tag.(string)
	}
	if sql == "" {
		sql = "select picture_name, tag from picture_tags where 1=1 " + CloudTags
	} else {
		sql = "select picture_name, tag from picture_tags where picture_name in (" + sql + ") " + CloudTags
	}

	log.Println(sql)
	rows, err := db.sess.Query("select name, filename from pictures where name in (select  picture_name from ("+sql+"))", tags...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name, filename, tag, desc string
		err = rows.Scan(&name, &filename)
		if err != nil {
			log.Println(err)
		}
		//predpokladame umisteni
		//pro fotky:  2015/02/web/_DSC1212.jpg
		//pro preview:  2015/02/thum/_DSC1212.jpg
		webname := strings.Replace(name, filename, "web/"+filename, 1)
		thumname := strings.Replace(name, filename, "thum/"+filename, 1)
		// description
		for _, meta := range options.descriptionTags {
			rows2, err := db.sess.Query("select tag from picture_tags where  picture_name = ? and meta = ? ", name, meta)
			if err != nil {
				log.Println(err)
			}
			defer rows2.Close()
			for rows2.Next() {
				err = rows2.Scan(&tag)
				if err != nil {
					log.Println(err)
				}
				desc = desc + ", " + tag
			}
		}
		sel.Images = append(sel.Images, Image{Image: "/pics/" + webname, Thumb: "/pics/" + thumname, Title: "", Description: strings.TrimPrefix(desc, ", ")})
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	var sel2 []Image
	// nahodne vybereme 100 obrazku
	// fmt.Println("Images ", len(sel.Images))
	for i, val := range rand.Perm(len(sel.Images)) {
		sel2 = append(sel2, sel.Images[val])
		if i > 100 {
			break
		}
	}
	// fmt.Println("Sel ", len(sel2))
	sel.Images = sel2
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	/*
		if len(sel.Images) <= 100 {
			sel.Tags = append(sel.Tags, Word{Text: "Gallery", Weight: 15, Link: options.link + "/show?" + strings.TrimPrefix(params, "&")})
		}
	*/
	rows, err = db.sess.Query(`
	select tag, cnt,
	       case 
		   when cnt between 1 and 10 then 1
		   when cnt between 11 and 20 then 2
		   when cnt between 21 and 50 then 3
		   when cnt between 51 and 100 then 4
		   when cnt between 101 and 200 then 5
		   when cnt between 201 and 400 then 6
		   when cnt between 401 and 800 then 7
		   when cnt > 800 then 8
		   end cnt_grp
    from ( select tag, count(distinct picture_name) cnt from (`+sql+") group by tag ) order by cnt_grp desc, tag", tags...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {

		var tag string
		var cnt, cntGrp int
		err = rows.Scan(&tag, &cnt, &cntGrp)
		if err != nil {
			log.Fatal(err)
		}
		sel.Tags = append(sel.Tags, Word{Text: tag, Weight: cntGrp, Count: cnt, Link: options.link + "/drilldown?" + strings.TrimPrefix(params+"&tag="+tag, "&")})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return sel
}
