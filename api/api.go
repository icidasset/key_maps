package api

import (
  "github.com/icidasset/key-maps/db"
  "github.com/pilu/traffic"
)

//
//  [API]
//  -> https://gobyexample.com/json
//  -> http://www.alexedwards.net/blog/golang-response-snippets
//  -> http://godoc.org/github.com/lib/pq
//

type Map struct {
  Id int `db:"id"`
}


//
//  [ROUTES]
//
func GetMaps(w traffic.ResponseWriter, r *traffic.Request) {
  maps := []Map{}
  db.Select(&maps, "SELECT * FROM maps")
}


func GetMap(w traffic.ResponseWriter, r *traffic.Request) {
  //
}
