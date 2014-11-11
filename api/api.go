package api

import (
	"database/sql"
	"github.com/pilu/traffic"
)

//
//  [API]
//  -> https://gobyexample.com/json
//  -> http://www.alexedwards.net/blog/golang-response-snippets
//  -> http://godoc.org/github.com/lib/pq
//

type Map struct {
	Id int
}


func GetMaps(w traffic.ResponseWriter, r *traffic.Request) {
	db, _ := sql.Open("postgres", "user=icidasset dbname=keymaps_development")
	defer db.Close()

	// db query
	rows, _ := db.Query("SELECT * FROM maps")
	defer rows.Close()

	// collect data
	for rows.Next() {
	}

	// output json
	w.WriteJSON(rows)
}


func GetMap(w traffic.ResponseWriter, r *traffic.Request) {
	// w.Write
}
