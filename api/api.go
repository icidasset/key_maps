package api

import (
  "encoding/json"
  "github.com/pilu/traffic"
)

//
//  [API]
//  -> https://gobyexample.com/json
//  -> http://www.alexedwards.net/blog/golang-response-snippets
//  -> http://godoc.org/github.com/lib/pq
//

type List struct {
  Id int
}


func renderJSON(w traffic.ResponseWriter, mj []byte) {
  w.Header().Set("Content-Type", "application/json")
  w.Write(mj)
}


func GetLists(w traffic.ResponseWriter, r *traffic.Request) {
  m := make(map[string]string)
  m["an"] = "example"
  m["a"] = "test"

  j, _ := json.Marshal(m)

  renderJSON(w, j)
}
