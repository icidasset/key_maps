package api

import (
  "fmt"
  "github.com/icidasset/key-maps/db"
  "github.com/pilu/traffic"
  "time"
)

//
//  [API]
//  -> https://gobyexample.com/json
//  -> http://www.alexedwards.net/blog/golang-response-snippets
//  -> http://godoc.org/github.com/lib/pq
//

type Map struct {
  Id int `db:"id"`
  Name string
  Structure string
  CreatedAt time.Time
  UpdatedAt time.Time
}


//
//  [ROUTES]
//
func GetMaps(w traffic.ResponseWriter, r *traffic.Request) {
  maps := []Map{}
  db.Select(&maps, "SELECT * FROM maps")
}


func GetMap(w traffic.ResponseWriter, r *traffic.Request) {
  m := Map{}
  db.Get(&m, "SELECT * FROM maps WHERE id = ?", r.Param("id"))

  w.WriteText( fmt.Sprintf("%#v\n", m) )
}
