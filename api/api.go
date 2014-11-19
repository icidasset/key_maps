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
  Id int
}


func BeforeFilter(w traffic.ResponseWriter, r *traffic.Request) {
  db.Open()
}


//
//  [ROUTES]
//
func GetMaps(w traffic.ResponseWriter, r *traffic.Request) {
  db.Close()
}


func GetMap(w traffic.ResponseWriter, r *traffic.Request) {
  db.Close()
}
