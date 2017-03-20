package main

import (
	"encoding/json"
	"log"
	"math"
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
	METADATA	text,
	NUM_METADATA	text
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
		stmt := `insert or replace into Pictures ( NAME, METADATA ) values( ?, ?)`
		meta, _ := json.Marshal(p.Metadata)
		_, err := db.sess.Exec(stmt, p.Name, meta)
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
	for i, tag := range tags {
		if i == 0 {
			sql = "select picture_name from picture_tags where tag = ?"
		} else if i < len(tags) {
			sql = "select picture_name from picture_tags where picture_name in (" + sql + ") and tag = ?"
		}
		params = params + "&tag=" + tag.(string)
	}
	if sql == "" {
		sql = "select picture_name, tag from picture_tags"
	} else {
		sql = "select picture_name, tag from picture_tags where picture_name in (" + sql + ")"
	}

	rows, err := db.sess.Query("select distinct picture_name from ("+sql+")", tags...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		sel.Images = append(sel.Images, Image{Image: "/pics/" + name, Thumb: "/pics/" + name})
		if len(sel.Images) > 100 {
			break
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if len(sel.Images) <= 100 {
		sel.Tags = append(sel.Tags, Word{Text: "Gallery", Weight: 15, Link: options.link + "/show?" + strings.TrimPrefix(params, "&")})
	}
	rows, err = db.sess.Query("select tag, count(distinct picture_name) cnt from ("+sql+") group by tag", tags...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tag string
		var cnt float64
		err = rows.Scan(&tag, &cnt)
		if err != nil {
			log.Fatal(err)
		}
		sel.Tags = append(sel.Tags, Word{Text: tag, Weight: int(math.Trunc(math.Log(cnt))), Link: options.link + "/drilldown?" + strings.TrimPrefix(params+"&tag="+tag, "&")})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return sel
}
